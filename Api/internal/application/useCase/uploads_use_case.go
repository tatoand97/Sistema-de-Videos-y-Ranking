package useCase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"api/internal/application/validations"
	"api/internal/domain"
	"api/internal/domain/entities"
	"api/internal/domain/interfaces"
	"api/internal/domain/requests"
	"api/internal/domain/responses"

	"github.com/google/uuid"
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
	publisher interfaces.MessagePublisher
	queue     string
}

func NewUploadsUseCase(videoRepo interfaces.VideoRepository, storage interfaces.VideoStorage, publisher interfaces.MessagePublisher, queue string) *UploadsUseCase {
	return &UploadsUseCase{videoRepo: videoRepo, storage: storage, publisher: publisher, queue: queue}
}

// sanitizeFileRe keeps only safe characters for object/file names.
var sanitizeFileRe = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func userIDFromCtx(ctx context.Context) (uint, error) {
	val := ctx.Value(UserIDContextKey)
	userID, ok := val.(uint)
	if !ok || userID == 0 {
		return 0, errors.New("userID missing in context")
	}
	return userID, nil
}

func readAllFromHeader(h *multipart.FileHeader) ([]byte, error) {
	f, err := h.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := make([]byte, h.Size)
	if _, err := io.ReadFull(f, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func sanitizeBaseName(name string) string {
	base := filepath.Base(name)
	if base == "." || base == ".." || strings.TrimSpace(base) == "" {
		base = "file"
	}
	base = strings.TrimSpace(base)
	base = strings.ReplaceAll(base, " ", "-")
	base = sanitizeFileRe.ReplaceAllString(base, "")
	if base == "" {
		base = "file"
	}
	return base
}

func buildObjectName(base string) string {
	return fmt.Sprintf("%s-%s", uuid.NewString(), base)
}

func contentTypeFromHeader(h *multipart.FileHeader) string {
	ct := h.Header.Get("Content-Type")
	if ct == "" {
		return "application/octet-stream"
	}
	return ct
}

func (uc *UploadsUseCase) saveFile(ctx context.Context, objectName string, data []byte, contentType string, size int64) (string, error) {
	reader := bytes.NewReader(data)
	return uc.storage.Save(ctx, objectName, reader, size, contentType)
}

func (uc *UploadsUseCase) persistVideo(ctx context.Context, userID uint, input UploadVideoInput, url string) (*entities.Video, error) {
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
	return video, nil
}

func (uc *UploadsUseCase) publishVideo(videoID uint) {
	// Publish the saved video ID for async processing.
	// Do not fail the upload if messaging fails.
	fmt.Printf("DEBUG: Checking publisher - publisher nil: %t, queue: '%s'\n", uc.publisher == nil, uc.queue)
	if uc.publisher == nil || strings.TrimSpace(uc.queue) == "" {
		if uc.publisher == nil {
			fmt.Printf("ERROR: Publisher is nil, cannot publish message for video ID: %d\n", videoID)
		} else {
			fmt.Printf("ERROR: Queue name is empty, cannot publish message for video ID: %d\n", videoID)
		}
		return
	}
	payload := struct {
		VideoID string `json:"videoId"`
	}{VideoID: fmt.Sprintf("%d", videoID)}

	b, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("ERROR: Failed to marshal message payload: %v\n", err)
		return
	}
	fmt.Printf("DEBUG: Publishing message for video ID: %d to queue: %s, payload: %s\n", videoID, uc.queue, string(b))
	if publishErr := uc.publisher.Publish(uc.queue, b); publishErr != nil {
		fmt.Printf("ERROR: Failed to publish message to queue %s: %v\n", uc.queue, publishErr)
		return
	}
	fmt.Printf("SUCCESS: Message published to queue %s for video ID: %d\n", uc.queue, videoID)
}

// UploadMultipart handles the classic multipart upload path.
func (uc *UploadsUseCase) UploadMultipart(ctx context.Context, input UploadVideoInput) (*UploadVideoOutput, error) {
	userID, err := userIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	fileBytes, err := readAllFromHeader(input.FileHeader)
	if err != nil {
		return nil, err
	}

	// Validate MP4
	if _, _, err := validations.CheckMP4(fileBytes); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalid, err)
	}

	// Build a unique object key to avoid overwriting files with same original name
	base := sanitizeBaseName(input.FileHeader.Filename)
	objectName := buildObjectName(base)
	contentType := contentTypeFromHeader(input.FileHeader)

	url, err := uc.saveFile(ctx, objectName, fileBytes, contentType, input.FileHeader.Size)
	if err != nil {
		return nil, err
	}

	video, err := uc.persistVideo(ctx, userID, input, url)
	if err != nil {
		return nil, err
	}

	uc.publishVideo(video.VideoID)

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
