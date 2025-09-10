package domain

import (
	"errors"
	"testing"
	"watermarking/internal/domain"

	"github.com/stretchr/testify/assert"
)

type TestVideoRepository struct{}

func (t *TestVideoRepository) FindByFilename(filename string) (*domain.Video, error) {
	if filename == "missing.mp4" {
		return nil, errors.New("video not found")
	}
	return &domain.Video{
		ID:       "123",
		Filename: filename,
		Status:   domain.StatusPending,
	}, nil
}

func (t *TestVideoRepository) UpdateStatus(id string, status domain.ProcessingStatus) error {
	if id == "invalid" {
		return errors.New("invalid video ID")
	}
	return nil
}

type TestStorageRepository struct{}

func (t *TestStorageRepository) Download(bucket, filename string) ([]byte, error) {
	if filename == "missing.mp4" {
		return nil, errors.New("file not found")
	}
	return []byte("video data"), nil
}

func (t *TestStorageRepository) Upload(bucket, filename string, data []byte) error {
	if bucket == "invalid-bucket" {
		return errors.New("bucket not found")
	}
	return nil
}

type TestProcessingService struct{}

func (t *TestProcessingService) TrimToMaxSeconds(data []byte, maxSeconds int) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	return []byte("watermarked data"), nil
}

type TestNotificationService struct{}

func (t *TestNotificationService) NotifyVideoProcessed(videoID, filename, bucketPath string) error {
	if videoID == "invalid" {
		return errors.New("invalid video ID")
	}
	return nil
}

func TestVideoRepository_Interface(t *testing.T) {
	var repo domain.VideoRepository = &TestVideoRepository{}
	assert.NotNil(t, repo)

	video, err := repo.FindByFilename("test.mp4")
	assert.NoError(t, err)
	assert.NotNil(t, video)
	assert.Equal(t, "123", video.ID)

	err = repo.UpdateStatus("123", domain.StatusCompleted)
	assert.NoError(t, err)
}

func TestStorageRepository_Interface(t *testing.T) {
	var repo domain.StorageRepository = &TestStorageRepository{}
	assert.NotNil(t, repo)

	data, err := repo.Download("bucket", "test.mp4")
	assert.NoError(t, err)
	assert.Equal(t, []byte("video data"), data)

	err = repo.Upload("bucket", "test.mp4", []byte("data"))
	assert.NoError(t, err)
}

func TestProcessingService_Interface(t *testing.T) {
	var service domain.VideoProcessingService = &TestProcessingService{}
	assert.NotNil(t, service)

	result, err := service.TrimToMaxSeconds([]byte("video data"), 30)
	assert.NoError(t, err)
	assert.Equal(t, []byte("watermarked data"), result)
}

func TestNotificationService_Interface(t *testing.T) {
	var service domain.NotificationService = &TestNotificationService{}
	assert.NotNil(t, service)

	err := service.NotifyVideoProcessed("123", "test.mp4", "bucket/test.mp4")
	assert.NoError(t, err)
}