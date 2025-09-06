package useCase

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"time"

	"io"
	"main_videork/internal/application/validations"
	"main_videork/internal/domain/entities"
	"main_videork/internal/domain/interfaces"
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

type UploadVideoUseCase struct {
	videoRepo interfaces.VideoRepository
	storage   interfaces.VideoStorage
}

func NewUploadVideoUseCase(videoRepo interfaces.VideoRepository, storage interfaces.VideoStorage) *UploadVideoUseCase {
	return &UploadVideoUseCase{videoRepo: videoRepo, storage: storage}
}

func (uc *UploadVideoUseCase) Execute(ctx context.Context, input UploadVideoInput) (*UploadVideoOutput, error) {
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

	// Leer el archivo en memoria para validación
	fileBytes := make([]byte, input.FileHeader.Size)
	if _, err := io.ReadFull(file, fileBytes); err != nil {
		return nil, err
	}

	// Validar MP4 (tamaño, resolución, brand)
	if _, _, err := validations.CheckMP4(fileBytes); err != nil {
		return nil, err
	}

	// Volver a crear un reader para el guardado (ya que file está consumido)
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
