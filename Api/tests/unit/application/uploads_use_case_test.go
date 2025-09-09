package application_test

import (
	"api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/entities"
	"api/internal/domain/requests"
	"api/tests/mocks"
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadsUseCase_UploadMultipart(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		input          useCase.UploadVideoInput
		mockRepo       *mocks.MockVideoRepository
		mockStorage    *mocks.MockVideoStorage
		mockPublisher  *mocks.MockMessagePublisher
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:   "successful upload",
			userID: 1,
			input: useCase.UploadVideoInput{
				Title:      "Test Video",
				FileHeader: createMockFileHeader("test.mp4", createValidMP4()),
				Status:     string(entities.StatusUploaded),
			},
			mockRepo: &mocks.MockVideoRepository{
				CreateFunc: func(ctx context.Context, video *entities.Video) error {
					video.VideoID = 123
					return nil
				},
			},
			mockStorage: &mocks.MockVideoStorage{
				SaveFunc: func(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
					return "https://example.com/video.mp4", nil
				},
			},
			mockPublisher: &mocks.MockMessagePublisher{},
			wantErr:       false,
		},
		{
			name:   "missing user ID",
			userID: 0,
			input: useCase.UploadVideoInput{
				Title:      "Test Video",
				FileHeader: createMockFileHeader("test.mp4", createValidMP4()),
				Status:     string(entities.StatusUploaded),
			},
			wantErr:        true,
			expectedErrMsg: "userID missing in context",
		},
		{
			name:   "invalid MP4 file",
			userID: 1,
			input: useCase.UploadVideoInput{
				Title:      "Test Video",
				FileHeader: createMockFileHeader("test.mp4", []byte("invalid content")),
				Status:     string(entities.StatusUploaded),
			},
			wantErr: true,
		},
		{
			name:   "storage save fails",
			userID: 1,
			input: useCase.UploadVideoInput{
				Title:      "Test Video",
				FileHeader: createMockFileHeader("test.mp4", createValidMP4()),
				Status:     string(entities.StatusUploaded),
			},
			mockStorage: &mocks.MockVideoStorage{
				SaveFunc: func(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
					return "", errors.New("storage error")
				},
			},
			wantErr: true,
		},
		{
			name:   "repository create fails",
			userID: 1,
			input: useCase.UploadVideoInput{
				Title:      "Test Video",
				FileHeader: createMockFileHeader("test.mp4", createValidMP4()),
				Status:     string(entities.StatusUploaded),
			},
			mockRepo: &mocks.MockVideoRepository{
				CreateFunc: func(ctx context.Context, video *entities.Video) error {
					return errors.New("database error")
				},
			},
			mockStorage: &mocks.MockVideoStorage{
				SaveFunc: func(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
					return "https://example.com/video.mp4", nil
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.userID > 0 {
				ctx = context.WithValue(ctx, useCase.UserIDContextKey, tt.userID)
			}

			usecase := useCase.NewUploadsUseCase(tt.mockRepo, tt.mockStorage, tt.mockPublisher, "test-queue")
			
			result, err := usecase.UploadMultipart(ctx, tt.input)
			
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.input.Title, result.Title)
				assert.NotZero(t, result.VideoID)
			}
		})
	}
}

func TestUploadsUseCase_CreatePostPolicy(t *testing.T) {
	tests := []struct {
		name        string
		req         requests.CreateUploadRequest
		mockStorage *mocks.MockVideoStorage
		wantErr     bool
	}{
		{
			name: "valid request",
			req: requests.CreateUploadRequest{
				Filename:  "test.mp4",
				MimeType:  "video/mp4",
				SizeBytes: 1024,
			},
			mockStorage: &mocks.MockVideoStorage{},
			wantErr:     false,
		},
		{
			name: "empty filename",
			req: requests.CreateUploadRequest{
				Filename:  "",
				MimeType:  "video/mp4",
				SizeBytes: 1024,
			},
			wantErr: true,
		},
		{
			name: "empty mime type",
			req: requests.CreateUploadRequest{
				Filename:  "test.mp4",
				MimeType:  "",
				SizeBytes: 1024,
			},
			wantErr: true,
		},
		{
			name: "negative size",
			req: requests.CreateUploadRequest{
				Filename:  "test.mp4",
				MimeType:  "video/mp4",
				SizeBytes: -1,
			},
			wantErr: true,
		},
		{
			name: "invalid checksum",
			req: requests.CreateUploadRequest{
				Filename:  "test.mp4",
				MimeType:  "video/mp4",
				SizeBytes: 1024,
				Checksum:  "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := useCase.NewUploadsUseCase(nil, tt.mockStorage, nil, "")
			
			result, err := usecase.CreatePostPolicy(context.Background(), tt.req)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestUploadsUseCase_DeleteUserVideoIfEligible(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		videoID  uint
		mockRepo *mocks.MockVideoRepository
		wantErr  bool
	}{
		{
			name:    "successful deletion",
			userID:  1,
			videoID: 123,
			mockRepo: &mocks.MockVideoRepository{
				GetByIDAndUserFunc: func(ctx context.Context, id, userID uint) (*entities.Video, error) {
					return &entities.Video{
						VideoID: id,
						UserID:  userID,
						Status:  string(entities.StatusUploaded),
					}, nil
				},
				DeleteFunc: func(ctx context.Context, id uint) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name:    "video not found",
			userID:  1,
			videoID: 999,
			mockRepo: &mocks.MockVideoRepository{
				GetByIDAndUserFunc: func(ctx context.Context, id, userID uint) (*entities.Video, error) {
					return nil, domain.ErrNotFound
				},
			},
			wantErr: true,
		},
		{
			name:    "video already processed",
			userID:  1,
			videoID: 123,
			mockRepo: &mocks.MockVideoRepository{
				GetByIDAndUserFunc: func(ctx context.Context, id, userID uint) (*entities.Video, error) {
					return &entities.Video{
						VideoID: id,
						UserID:  userID,
						Status:  string(entities.StatusProcessed),
					}, nil
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := useCase.NewUploadsUseCase(tt.mockRepo, nil, nil, "")
			
			err := usecase.DeleteUserVideoIfEligible(context.Background(), tt.userID, tt.videoID)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper functions
func createMockFileHeader(filename string, content []byte) *multipart.FileHeader {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
	h.Set("Content-Type", "video/mp4")
	
	part, _ := writer.CreatePart(h)
	part.Write(content)
	writer.Close()
	
	reader := multipart.NewReader(body, writer.Boundary())
	form, _ := reader.ReadForm(int64(len(content)) + 1024)
	
	return form.File["file"][0]
}

func createValidMP4() []byte {
	// Minimal valid MP4 header
	return []byte{
		0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, // ftyp box
		0x69, 0x73, 0x6f, 0x6d, 0x00, 0x00, 0x02, 0x00,
		0x69, 0x73, 0x6f, 0x6d, 0x69, 0x73, 0x6f, 0x32,
		0x61, 0x76, 0x63, 0x31, 0x6d, 0x70, 0x34, 0x31,
	}
}