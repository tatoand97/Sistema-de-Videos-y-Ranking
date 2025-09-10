package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test implementations for interfaces
type TestVideoRepository struct{}

func (t *TestVideoRepository) FindByFilename(filename string) (*Video, error) {
	if filename == "nonexistent.mp4" {
		return nil, errors.New("video not found")
	}
	return &Video{
		ID:       "123",
		Filename: filename,
		Status:   StatusPending,
	}, nil
}

func (t *TestVideoRepository) UpdateStatus(videoID string, status ProcessingStatus) error {
	if videoID == "invalid" {
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

func (t *TestProcessingService) RemoveAudio(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	return []byte("processed data"), nil
}

type TestNotificationService struct{}

func (t *TestNotificationService) NotifyVideoProcessed(videoID, filename, bucketPath string) error {
	if videoID == "invalid" {
		return errors.New("invalid video ID")
	}
	return nil
}

func (t *TestNotificationService) NotifyProcessingComplete(videoID string, success bool) error {
	if videoID == "invalid" {
		return errors.New("invalid video ID")
	}
	return nil
}

func TestVideoRepository_Interface(t *testing.T) {
	var repo VideoRepository = &TestVideoRepository{}
	assert.NotNil(t, repo)

	// Test FindByFilename success
	video, err := repo.FindByFilename("test.mp4")
	assert.NoError(t, err)
	assert.NotNil(t, video)
	assert.Equal(t, "123", video.ID)

	// Test FindByFilename error
	video, err = repo.FindByFilename("nonexistent.mp4")
	assert.Error(t, err)
	assert.Nil(t, video)

	// Test UpdateStatus success
	err = repo.UpdateStatus("123", StatusCompleted)
	assert.NoError(t, err)

	// Test UpdateStatus error
	err = repo.UpdateStatus("invalid", StatusCompleted)
	assert.Error(t, err)
}

func TestStorageRepository_Interface(t *testing.T) {
	var repo StorageRepository = &TestStorageRepository{}
	assert.NotNil(t, repo)

	// Test Download success
	data, err := repo.Download("bucket", "test.mp4")
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, []byte("video data"), data)

	// Test Download error
	data, err = repo.Download("bucket", "missing.mp4")
	assert.Error(t, err)
	assert.Nil(t, data)

	// Test Upload success
	err = repo.Upload("bucket", "test.mp4", []byte("data"))
	assert.NoError(t, err)

	// Test Upload error
	err = repo.Upload("invalid-bucket", "test.mp4", []byte("data"))
	assert.Error(t, err)
}

func TestProcessingService_Interface(t *testing.T) {
	var service VideoProcessingService = &TestProcessingService{}
	assert.NotNil(t, service)

	// Test RemoveAudio success
	result, err := service.RemoveAudio([]byte("video data"))
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, []byte("processed data"), result)

	// Test RemoveAudio error
	result, err = service.RemoveAudio([]byte{})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNotificationService_Interface(t *testing.T) {
	var service NotificationService = &TestNotificationService{}
	assert.NotNil(t, service)

	// Test NotifyVideoProcessed success
	err := service.NotifyVideoProcessed("123", "test.mp4", "bucket/test.mp4")
	assert.NoError(t, err)

	// Test NotifyVideoProcessed error
	err = service.NotifyVideoProcessed("invalid", "test.mp4", "bucket/test.mp4")
	assert.Error(t, err)

	// Test NotifyProcessingComplete success
	err = service.NotifyProcessingComplete("123", true)
	assert.NoError(t, err)

	// Test NotifyProcessingComplete error
	err = service.NotifyProcessingComplete("invalid", true)
	assert.Error(t, err)
}