package domain

import (
	"statesmachine/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test interfaces exist and can be implemented
type TestVideoRepository struct{}

func (t *TestVideoRepository) FindByID(id uint) (*domain.Video, error) {
	return &domain.Video{ID: id, OriginalFile: "test.mp4", Status: "UPLOADED"}, nil
}

func (t *TestVideoRepository) UpdateStatus(id uint, status domain.VideoStatus) error {
	return nil
}

func (t *TestVideoRepository) UpdateStatusAndProcessedFile(id uint, status domain.VideoStatus, processedFile string) error {
	return nil
}

type TestMessagePublisher struct{}

func (t *TestMessagePublisher) PublishMessage(queue string, message []byte) error {
	return nil
}

func TestVideoRepository_Interface(t *testing.T) {
	var repo domain.VideoRepository = &TestVideoRepository{}
	assert.NotNil(t, repo)

	// Test FindByID
	video, err := repo.FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, video)
	assert.Equal(t, uint(1), video.ID)

	// Test UpdateStatus
	err = repo.UpdateStatus(1, domain.StatusTrimming)
	assert.NoError(t, err)

	// Test UpdateStatusAndProcessedFile
	err = repo.UpdateStatusAndProcessedFile(1, domain.StatusProcessed, "final.mp4")
	assert.NoError(t, err)
}

func TestMessagePublisher_Interface(t *testing.T) {
	var publisher domain.MessagePublisher = &TestMessagePublisher{}
	assert.NotNil(t, publisher)

	// Test PublishMessage
	err := publisher.PublishMessage("test_queue", []byte("test message"))
	assert.NoError(t, err)
}