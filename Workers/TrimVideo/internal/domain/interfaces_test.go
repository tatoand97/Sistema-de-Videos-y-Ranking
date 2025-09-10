package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestVideoRepository struct{}

func (t *TestVideoRepository) FindByFilename(filename string) (*Video, error) {
	if filename == "missing.mp4" {
		return nil, errors.New("video not found")
	}
	return &Video{
		ID:       "123",
		Filename: filename,
		Status:   StatusPending,
	}, nil
}

func (t *TestVideoRepository) UpdateStatus(id string, status ProcessingStatus) error {
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
	return []byte("trimmed data"), nil
}

type TestNotificationService struct{}

func (t *TestNotificationService) NotifyVideoProcessed(videoID, filename, bucketPath string) error {
	if videoID == "invalid" {
		return errors.New("invalid video ID")
	}
	return nil
}

func TestVideoRepository_Interface(t *testing.T) {
	var repo VideoRepository = &TestVideoRepository{}
	assert.NotNil(t, repo)

	video, err := repo.FindByFilename("test.mp4")
	assert.NoError(t, err)
	assert.NotNil(t, video)
	assert.Equal(t, "123", video.ID)

	video, err = repo.FindByFilename("missing.mp4")
	assert.Error(t, err)
	assert.Nil(t, video)

	err = repo.UpdateStatus("123", StatusCompleted)
	assert.NoError(t, err)

	err = repo.UpdateStatus("invalid", StatusCompleted)
	assert.Error(t, err)
}

func TestStorageRepository_Interface(t *testing.T) {
	var repo StorageRepository = &TestStorageRepository{}
	assert.NotNil(t, repo)

	data, err := repo.Download("bucket", "test.mp4")
	assert.NoError(t, err)
	assert.Equal(t, []byte("video data"), data)

	data, err = repo.Download("bucket", "missing.mp4")
	assert.Error(t, err)
	assert.Nil(t, data)

	err = repo.Upload("bucket", "test.mp4", []byte("data"))
	assert.NoError(t, err)

	err = repo.Upload("invalid-bucket", "test.mp4", []byte("data"))
	assert.Error(t, err)
}

func TestProcessingService_Interface(t *testing.T) {
	var service VideoProcessingService = &TestProcessingService{}
	assert.NotNil(t, service)

	result, err := service.TrimToMaxSeconds([]byte("video data"), 30)
	assert.NoError(t, err)
	assert.Equal(t, []byte("trimmed data"), result)

	result, err = service.TrimToMaxSeconds([]byte{}, 30)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNotificationService_Interface(t *testing.T) {
	var service NotificationService = &TestNotificationService{}
	assert.NotNil(t, service)

	err := service.NotifyVideoProcessed("123", "test.mp4", "bucket/test.mp4")
	assert.NoError(t, err)

	err = service.NotifyVideoProcessed("invalid", "test.mp4", "bucket/test.mp4")
	assert.Error(t, err)
}