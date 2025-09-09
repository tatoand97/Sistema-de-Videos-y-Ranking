package useCase

import (
	"api/internal/domain/entities"
	"api/internal/domain/requests"
	"api/tests/mocks"
	"api/tests/testdata"
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUploadsUseCase(t *testing.T) {
	repo := &mocks.MockVideoRepository{}
	storage := &mocks.MockVideoStorage{}
	publisher := &mocks.MockMessagePublisher{}
	queueName := "test-queue"
	
	usecase := NewUploadsUseCase(repo, storage, publisher, queueName)
	
	assert.NotNil(t, usecase)
}

func TestUploadsUseCase_UploadMultipart_Success(t *testing.T) {
	mockRepo := &mocks.MockVideoRepository{
		CreateFunc: func(ctx context.Context, video *entities.Video) error {
			video.VideoID = 123
			return nil
		},
	}
	mockStorage := &mocks.MockVideoStorage{
		SaveFunc: func(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
			return "https://example.com/video.mp4", nil
		},
	}
	mockPublisher := &mocks.MockMessagePublisher{}
	
	usecase := NewUploadsUseCase(mockRepo, mockStorage, mockPublisher, "test-queue")
	
	ctx := context.WithValue(context.Background(), UserIDContextKey, uint(1))
	input := UploadVideoInput{
		Title:      "Test Video",
		FileHeader: createMockFileHeader("test.mp4", testdata.CreateValidMP4()),
		Status:     string(entities.StatusUploaded),
	}
	
	result, err := usecase.UploadMultipart(ctx, input)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Video", result.Title)
	assert.Equal(t, uint(123), result.VideoID)
}

func TestUploadsUseCase_CreatePostPolicy_Success(t *testing.T) {
	mockStorage := &mocks.MockVideoStorage{}
	
	usecase := NewUploadsUseCase(nil, mockStorage, nil, "")
	
	req := requests.CreateUploadRequest{
		Filename:  "test.mp4",
		MimeType:  "video/mp4",
		SizeBytes: 1024,
	}
	
	result, err := usecase.CreatePostPolicy(context.Background(), req)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUploadsUseCase_CreatePostPolicy_ValidationErrors(t *testing.T) {
	usecase := NewUploadsUseCase(nil, &mocks.MockVideoStorage{}, nil, "")
	
	tests := []struct {
		name string
		req  requests.CreateUploadRequest
	}{
		{
			name: "empty filename",
			req: requests.CreateUploadRequest{
				Filename:  "",
				MimeType:  "video/mp4",
				SizeBytes: 1024,
			},
		},
		{
			name: "empty mime type",
			req: requests.CreateUploadRequest{
				Filename:  "test.mp4",
				MimeType:  "",
				SizeBytes: 1024,
			},
		},
		{
			name: "negative size",
			req: requests.CreateUploadRequest{
				Filename:  "test.mp4",
				MimeType:  "video/mp4",
				SizeBytes: -1,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := usecase.CreatePostPolicy(context.Background(), tt.req)
			
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestUploadsUseCase_ListUserVideos(t *testing.T) {
	expectedVideos := []*entities.Video{
		{VideoID: 1, Title: "Video 1", UserID: 123},
		{VideoID: 2, Title: "Video 2", UserID: 123},
	}
	
	mockRepo := &mocks.MockVideoRepository{
		ListByUserFunc: func(ctx context.Context, userID uint) ([]*entities.Video, error) {
			return expectedVideos, nil
		},
	}
	
	usecase := NewUploadsUseCase(mockRepo, nil, nil, "")
	
	videos, err := usecase.ListUserVideos(context.Background(), 123)
	
	assert.NoError(t, err)
	assert.Equal(t, expectedVideos, videos)
}

func TestUploadsUseCase_GetUserVideoByID(t *testing.T) {
	expectedVideo := &entities.Video{VideoID: 1, Title: "Test Video", UserID: 123}
	
	mockRepo := &mocks.MockVideoRepository{
		GetByIDAndUserFunc: func(ctx context.Context, id, userID uint) (*entities.Video, error) {
			return expectedVideo, nil
		},
	}
	
	usecase := NewUploadsUseCase(mockRepo, nil, nil, "")
	
	video, err := usecase.GetUserVideoByID(context.Background(), 1, 123)
	
	assert.NoError(t, err)
	assert.Equal(t, expectedVideo, video)
}

func TestUploadsUseCase_DeleteUserVideoIfEligible(t *testing.T) {
	mockRepo := &mocks.MockVideoRepository{
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
	}
	
	usecase := NewUploadsUseCase(mockRepo, nil, nil, "")
	
	err := usecase.DeleteUserVideoIfEligible(context.Background(), 123, 1)
	
	assert.NoError(t, err)
}

func TestUploadsUseCase_DeleteUserVideoIfEligible_AlreadyProcessed(t *testing.T) {
	mockRepo := &mocks.MockVideoRepository{
		GetByIDAndUserFunc: func(ctx context.Context, id, userID uint) (*entities.Video, error) {
			return &entities.Video{
				VideoID: id,
				UserID:  userID,
				Status:  string(entities.StatusProcessed),
			}, nil
		},
	}
	
	usecase := NewUploadsUseCase(mockRepo, nil, nil, "")
	
	err := usecase.DeleteUserVideoIfEligible(context.Background(), 123, 1)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid input")
}

// Helper function
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