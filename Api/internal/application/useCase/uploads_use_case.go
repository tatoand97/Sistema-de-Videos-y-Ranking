package useCase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"regexp"
	"strings"
	"time"

	"api/internal/application/validations"
	"api/internal/domain"
	"api/internal/domain/entities"
	"api/internal/domain/interfaces"
	"api/internal/domain/requests"
	"api/internal/domain/responses"
)

type UploadVideoInput struct {
	Title      string
	FileHeader *multipart.FileHeader
	Status     string
}

type contextKey string

// UserIDContextKey is the key used to store the authenticated user's ID in the context.
const UserIDContextKey contextKey = "userID"

type UploadVideoOutput struct {
	VideoID      uint
	Title        string
	OriginalFile string
	UploadedAt   time.Time
}

type UploadsUseCase struct {
	videoRepo interfaces.VideoRepository
	storage   interfaces.VideoStorage
}

func NewUploadsUseCase(videoRepo interfaces.VideoRepository, storage interfaces.VideoStorage) *UploadsUseCase {
	return &UploadsUseCase{videoRepo: videoRepo, storage: storage}
}

// UploadMultipart handles the classic multipart upload path.
func (uc *UploadsUseCase) UploadMultipart(ctx context.Context, input UploadVideoInput) (*UploadVideoOutput, error) {
	val := ctx.Value(UserIDContextKey)
	userID, ok := val.(uint)
	if !ok || userID == 0 {
		return nil, errors.New("userID missing in context")
	}

	file, err := input.FileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read into memory for validation
	fileBytes := make([]byte, input.FileHeader.Size)
	if _, err := io.ReadFull(file, fileBytes); err != nil {
		return nil, err
	}

	// Validate MP4
	if _, _, err := validations.CheckMP4(fileBytes); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalid, err)
	}

	// Recreate a reader for saving
	reader := bytes.NewReader(fileBytes)

	objectName := input.FileHeader.Filename
	contentType := input.FileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	url, err := uc.storage.Save(ctx, objectName, reader, input.FileHeader.Size, contentType)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	video := &entities.Video{
		UserID:       userID,
		Title:        input.Title,
		OriginalFile: url,
		Status:       input.Status,
		UploadedAt:   now,
	}
	if err := uc.videoRepo.Create(ctx, video); err != nil {
		return nil, err
	}

	return &UploadVideoOutput{
		VideoID:      video.VideoID,
		Title:        video.Title,
		OriginalFile: video.OriginalFile,
		UploadedAt:   video.UploadedAt,
	}, nil
}

var sha256HexRe = regexp.MustCompile(`^[A-Fa-f0-9]{64}$`)

// CreatePostPolicy validates input and delegates to storage to build a POST policy.
func (uc *UploadsUseCase) CreatePostPolicy(ctx context.Context, req requests.CreateUploadRequest) (*responses.CreateUploadResponsePostPolicy, error) {
	if strings.TrimSpace(req.Filename) == "" || strings.TrimSpace(req.MimeType) == "" {
		return nil, fmt.Errorf("%w: filename and mimeType are required", domain.ErrInvalid)
	}
	if req.SizeBytes < 0 {
		return nil, fmt.Errorf("%w: sizeBytes must be >= 0", domain.ErrInvalid)
	}
	if req.Checksum != "" && !sha256HexRe.MatchString(req.Checksum) {
		return nil, fmt.Errorf("%w: checksum must be SHA-256 hex", domain.ErrInvalid)
	}
	return uc.storage.PresignedPostPolicy(ctx, req)
}

// ListUserVideos returns all videos belonging to a user.
func (uc *UploadsUseCase) ListUserVideos(ctx context.Context, userID uint) ([]*entities.Video, error) {
	return uc.videoRepo.ListByUser(ctx, userID)
}

// GetUserVideoByID returns a user's video by id, enforcing ownership.
func (uc *UploadsUseCase) GetUserVideoByID(ctx context.Context, userID, videoID uint) (*entities.Video, error) {
	return uc.videoRepo.GetByIDAndUser(ctx, videoID, userID)
}

// DeleteUserVideoIfEligible deletes a user's own video if it meets business rules.
// Rules: owner-only, and cannot delete if the video is already processed (published for voting).
func (uc *UploadsUseCase) DeleteUserVideoIfEligible(ctx context.Context, userID, videoID uint) error {
	// Ensure ownership and existence
	v, err := uc.videoRepo.GetByIDAndUser(ctx, videoID, userID)
	if err != nil {
		return err
	}
	// Eligibility: not allowed if processed (published for voting)
	if v.Status == string(entities.StatusProcessed) {
		return domain.ErrInvalid
	}
	// Proceed to delete
	if err := uc.videoRepo.Delete(ctx, videoID); err != nil {
		return err
	}
	return nil
}
