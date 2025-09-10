package adapters

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
	"statesmachine/internal/adapters"

	"github.com/stretchr/testify/assert"
)

func TestMessageHandler_HandleMessage_EditCompleted(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	processedMsg := adapters.VideoProcessedMessage{
		VideoID:    "123",
		Filename:   "edited.mp4",
		BucketPath: "edit/edited.mp4",
		Status:     "completed",
	}

	msgBytes, _ := json.Marshal(processedMsg)
	mockUC.On("HandleEditCompleted", "123", "edited.mp4").Return(nil)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_AudioRemovalCompleted(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	processedMsg := adapters.VideoProcessedMessage{
		VideoID:    "123",
		Filename:   "no-audio.mp4",
		BucketPath: "audio-removal/no-audio.mp4",
		Status:     "completed",
	}

	msgBytes, _ := json.Marshal(processedMsg)
	mockUC.On("HandleAudioRemovalCompleted", "123", "no-audio.mp4").Return(nil)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_WatermarkingCompleted(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	processedMsg := adapters.VideoProcessedMessage{
		VideoID:    "123",
		Filename:   "watermarked.mp4",
		BucketPath: "watermarking/watermarked.mp4",
		Status:     "completed",
	}

	msgBytes, _ := json.Marshal(processedMsg)
	mockUC.On("HandleWatermarkingCompleted", "123", "watermarked.mp4").Return(nil)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_GossipCompleted(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	processedMsg := adapters.VideoProcessedMessage{
		VideoID:    "123",
		Filename:   "final.mp4",
		BucketPath: "processed-videos/final.mp4",
		Status:     "completed",
	}

	msgBytes, _ := json.Marshal(processedMsg)
	mockUC.On("HandleGossipOpenCloseCompleted", "123", "final.mp4").Return(nil)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_ProcessedInvalidVideoID(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	processedMsg := adapters.VideoProcessedMessage{
		VideoID:    "invalid",
		Filename:   "test.mp4",
		BucketPath: "trim/test.mp4",
		Status:     "completed",
	}

	msgBytes, _ := json.Marshal(processedMsg)
	mockUC.On("HandleTrimCompleted", "invalid", "test.mp4").Return(errors.New("invalid video ID format"))

	err := handler.HandleMessage(msgBytes)

	assert.Error(t, err)
	assert.IsType(t, &adapters.NonRetryableError{}, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_MaxRetriesExceeded(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	msg := adapters.VideoMessage{
		VideoID:    "123",
		Filename:   "test.mp4",
		RetryCount: 5,
		MaxRetries: 3,
	}

	msgBytes, _ := json.Marshal(msg)

	err := handler.HandleMessage(msgBytes)

	assert.Error(t, err)
	assert.IsType(t, &adapters.NonRetryableError{}, err)
	assert.Contains(t, err.Error(), "Max retries exceeded")
}

func TestMessageHandler_HandleMessage_RetryDelayNotMet(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	msg := adapters.VideoMessage{
		VideoID:    "123",
		Filename:   "test.mp4",
		RetryCount: 1,
		MaxRetries: 3,
		LastRetry:  time.Now().Unix() - 60, // 1 minute ago
	}

	msgBytes, _ := json.Marshal(msg)
	mockUC.On("GetRetryDelayMinutes").Return(5) // 5 minutes required

	err := handler.HandleMessage(msgBytes)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retry delay not met")
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_RetryDelayMet(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	msg := adapters.VideoMessage{
		VideoID:    "123",
		Filename:   "test.mp4",
		RetryCount: 1,
		MaxRetries: 3,
		LastRetry:  time.Now().Unix() - 600, // 10 minutes ago
	}

	msgBytes, _ := json.Marshal(msg)
	mockUC.On("GetRetryDelayMinutes").Return(5) // 5 minutes required
	mockUC.On("Execute", "123").Return(nil)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_ExecuteInvalidVideoID(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	msg := adapters.VideoMessage{
		VideoID:  "invalid",
		Filename: "test.mp4",
	}

	msgBytes, _ := json.Marshal(msg)
	mockUC.On("Execute", "invalid").Return(errors.New("invalid video ID format"))

	err := handler.HandleMessage(msgBytes)

	assert.Error(t, err)
	assert.IsType(t, &adapters.NonRetryableError{}, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_ExecuteOtherError(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	msg := adapters.VideoMessage{
		VideoID:  "123",
		Filename: "test.mp4",
	}

	msgBytes, _ := json.Marshal(msg)
	mockUC.On("Execute", "123").Return(errors.New("database connection error"))

	err := handler.HandleMessage(msgBytes)

	assert.Error(t, err)
	assert.NotEqual(t, &adapters.NonRetryableError{}, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_ProcessedUnknownBucketPath(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	processedMsg := adapters.VideoProcessedMessage{
		VideoID:    "123",
		Filename:   "test.mp4",
		BucketPath: "unknown/test.mp4",
		Status:     "completed",
	}

	msgBytes, _ := json.Marshal(processedMsg)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
}