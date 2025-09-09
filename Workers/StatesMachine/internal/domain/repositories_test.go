package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test interfaces exist and can be implemented
type TestVideoRepository struct{}

func (t *TestVideoRepository) FindByID(id uint) (*Video, error) {
	return &Video{ID: id, OriginalFile: "test.mp4", Status: "UPLOADED"}, nil
}

func (t *TestVideoRepository) UpdateStatus(id uint, status VideoStatus) error {
	return nil
}

func (t *TestVideoRepository) UpdateStatusAndProcessedFile(id uint, status VideoStatus, processedFile string) error {
	return nil
}

type TestMessagePublisher struct{}

func (t *TestMessagePublisher) PublishMessage(queue string, message []byte) error {
	return nil
}

func TestVideoRepository_Interface(t *testing.T) {
	var repo VideoRepository = &TestVideoRepository{}
	assert.NotNil(t, repo)

	// Test FindByID
	video, err := repo.FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, video)
	assert.Equal(t, uint(1), video.ID)

	// Test UpdateStatus
	err = repo.UpdateStatus(1, StatusTrimming)
	assert.NoError(t, err)

	// Test UpdateStatusAndProcessedFile
	err = repo.UpdateStatusAndProcessedFile(1, StatusProcessed, "final.mp4")
	assert.NoError(t, err)
}

func TestMessagePublisher_Interface(t *testing.T) {
	var publisher MessagePublisher = &TestMessagePublisher{}
	assert.NotNil(t, publisher)

	// Test PublishMessage
	err := publisher.PublishMessage("test_queue", []byte("test message"))
	assert.NoError(t, err)
}

func TestVideoRepository_Methods(t *testing.T) {
	repo := &TestVideoRepository{}

	// Test method signatures
	video, err := repo.FindByID(123)
	assert.NoError(t, err)
	assert.Equal(t, uint(123), video.ID)
	assert.Equal(t, "test.mp4", video.OriginalFile)
	assert.Equal(t, "UPLOADED", video.Status)

	err = repo.UpdateStatus(123, StatusProcessed)
	assert.NoError(t, err)

	err = repo.UpdateStatusAndProcessedFile(123, StatusProcessed, "processed.mp4")
	assert.NoError(t, err)
}

func TestMessagePublisher_Methods(t *testing.T) {
	publisher := &TestMessagePublisher{}

	// Test method signature
	err := publisher.PublishMessage("video_queue", []byte(`{"video_id": "123"}`))
	assert.NoError(t, err)
}

func TestInterfaces_Compatibility(t *testing.T) {
	// Test that our test implementations satisfy the interfaces
	var videoRepo VideoRepository
	var msgPublisher MessagePublisher

	videoRepo = &TestVideoRepository{}
	msgPublisher = &TestMessagePublisher{}

	assert.NotNil(t, videoRepo)
	assert.NotNil(t, msgPublisher)

	// Test interface methods work
	video, err := videoRepo.FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, video)

	err = msgPublisher.PublishMessage("test", []byte("test"))
	assert.NoError(t, err)
}