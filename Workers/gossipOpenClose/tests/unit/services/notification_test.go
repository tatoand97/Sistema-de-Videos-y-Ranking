package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessagePublisher struct {
	mock.Mock
}

func (m *MockMessagePublisher) PublishMessage(queue string, message interface{}) error {
	args := m.Called(queue, message)
	return args.Error(0)
}

func (m *MockMessagePublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewNotificationService(t *testing.T) {
	publisher := &MockMessagePublisher{}
	stateQueue := "state-queue"
	
	service := NewNotificationService(publisher, stateQueue)
	
	assert.NotNil(t, service)
	assert.Equal(t, stateQueue, service.stateQueue)
	assert.Equal(t, publisher, service.publisher)
}

func TestNotificationService_NotifyVideoProcessed_Success(t *testing.T) {
	publisher := &MockMessagePublisher{}
	service := NewNotificationService(publisher, "state-queue")
	
	videoID := "video-123"
	filename := "test.mp4"
	bucketPath := "processed/test.mp4"
	
	expectedMsg := VideoProcessedMessage{
		VideoID:    videoID,
		Filename:   filename,
		BucketPath: bucketPath,
		Status:     "completed",
	}
	
	publisher.On("PublishMessage", "state-queue", expectedMsg).Return(nil)
	
	err := service.NotifyVideoProcessed(videoID, filename, bucketPath)
	
	assert.NoError(t, err)
	publisher.AssertExpectations(t)
}

func TestNotificationService_NotifyVideoProcessed_Error(t *testing.T) {
	publisher := &MockMessagePublisher{}
	service := NewNotificationService(publisher, "state-queue")
	
	videoID := "video-123"
	filename := "test.mp4"
	bucketPath := "processed/test.mp4"
	expectedError := errors.New("publish failed")
	
	expectedMsg := VideoProcessedMessage{
		VideoID:    videoID,
		Filename:   filename,
		BucketPath: bucketPath,
		Status:     "completed",
	}
	
	publisher.On("PublishMessage", "state-queue", expectedMsg).Return(expectedError)
	
	err := service.NotifyVideoProcessed(videoID, filename, bucketPath)
	
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	publisher.AssertExpectations(t)
}

func TestVideoProcessedMessage_Structure(t *testing.T) {
	msg := VideoProcessedMessage{
		VideoID:    "test-id",
		Filename:   "test.mp4",
		BucketPath: "bucket/path",
		Status:     "completed",
	}
	
	assert.Equal(t, "test-id", msg.VideoID)
	assert.Equal(t, "test.mp4", msg.Filename)
	assert.Equal(t, "bucket/path", msg.BucketPath)
	assert.Equal(t, "completed", msg.Status)
}