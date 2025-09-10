package services_test

import (
	"audioremoval/internal/application/services"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockMessagePublisher struct {
	PublishMessageFunc func(queue string, message interface{}) error
	PublishedMessages  []PublishedMessage
	ShouldFail         bool
}

type PublishedMessage struct {
	Queue   string
	Message interface{}
}

func (m *MockMessagePublisher) PublishMessage(queue string, message interface{}) error {
	m.PublishedMessages = append(m.PublishedMessages, PublishedMessage{
		Queue:   queue,
		Message: message,
	})

	if m.PublishMessageFunc != nil {
		return m.PublishMessageFunc(queue, message)
	}

	if m.ShouldFail {
		return errors.New("publish failed")
	}

	return nil
}

func TestNewNotificationService(t *testing.T) {
	publisher := &MockMessagePublisher{}
	service := services.NewNotificationService(publisher, "test-queue")

	assert.NotNil(t, service)
}

func TestNotificationService_NotifyVideoProcessed_Success(t *testing.T) {
	// Arrange
	publisher := &MockMessagePublisher{}
	service := services.NewNotificationService(publisher, "state-queue")

	// Act
	err := service.NotifyVideoProcessed("video-123", "test.mp4", "processed-bucket/test.mp4")

	// Assert
	require.NoError(t, err)
	assert.Len(t, publisher.PublishedMessages, 1)
	assert.Equal(t, "state-queue", publisher.PublishedMessages[0].Queue)

	msg, ok := publisher.PublishedMessages[0].Message.(services.VideoProcessedMessage)
	require.True(t, ok)
	assert.Equal(t, "video-123", msg.VideoID)
	assert.Equal(t, "test.mp4", msg.Filename)
	assert.Equal(t, "processed-bucket/test.mp4", msg.BucketPath)
	assert.Equal(t, "completed", msg.Status)
}

func TestNotificationService_NotifyVideoProcessed_PublishFails(t *testing.T) {
	// Arrange
	publisher := &MockMessagePublisher{
		ShouldFail: true,
	}
	service := services.NewNotificationService(publisher, "state-queue")

	// Act
	err := service.NotifyVideoProcessed("video-123", "test.mp4", "processed-bucket/test.mp4")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "publish failed")
}

func TestNotificationService_NotifyProcessingComplete(t *testing.T) {
	// Arrange
	publisher := &MockMessagePublisher{}
	service := services.NewNotificationService(publisher, "state-queue")

	// Act
	err := service.NotifyProcessingComplete("video-123", true)

	// Assert
	require.NoError(t, err)
}

func TestVideoProcessedMessage_Structure(t *testing.T) {
	msg := services.VideoProcessedMessage{
		VideoID:    "video-123",
		Filename:   "test.mp4",
		BucketPath: "bucket/path",
		Status:     "completed",
	}

	assert.Equal(t, "video-123", msg.VideoID)
	assert.Equal(t, "test.mp4", msg.Filename)
	assert.Equal(t, "bucket/path", msg.BucketPath)
	assert.Equal(t, "completed", msg.Status)
}