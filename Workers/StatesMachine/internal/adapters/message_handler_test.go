package adapters

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrchestrateUseCase struct {
	mock.Mock
}

func (m *MockOrchestrateUseCase) Execute(videoID string) error {
	args := m.Called(videoID)
	return args.Error(0)
}

func (m *MockOrchestrateUseCase) HandleTrimCompleted(videoID, filename string) error {
	args := m.Called(videoID, filename)
	return args.Error(0)
}

func (m *MockOrchestrateUseCase) HandleEditCompleted(videoID, filename string) error {
	args := m.Called(videoID, filename)
	return args.Error(0)
}

func (m *MockOrchestrateUseCase) HandleAudioRemovalCompleted(videoID, filename string) error {
	args := m.Called(videoID, filename)
	return args.Error(0)
}

func (m *MockOrchestrateUseCase) HandleWatermarkingCompleted(videoID, filename string) error {
	args := m.Called(videoID, filename)
	return args.Error(0)
}

func (m *MockOrchestrateUseCase) HandleGossipOpenCloseCompleted(videoID, filename string) error {
	args := m.Called(videoID, filename)
	return args.Error(0)
}

func (m *MockOrchestrateUseCase) GetRetryDelayMinutes() int {
	args := m.Called()
	return args.Int(0)
}

func TestNonRetryableError_Error(t *testing.T) {
	originalErr := errors.New("original error")
	nonRetryable := &NonRetryableError{
		OriginalError: originalErr,
		Message:       "test message",
	}

	expected := "non-retryable error: test message - original error"
	assert.Equal(t, expected, nonRetryable.Error())
}

func TestNonRetryableError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	nonRetryable := &NonRetryableError{
		OriginalError: originalErr,
		Message:       "test message",
	}

	assert.Equal(t, originalErr, nonRetryable.Unwrap())
}

func TestIsNonRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "NonRetryableError type",
			err:      &NonRetryableError{OriginalError: errors.New("test"), Message: "test"},
			expected: true,
		},
		{
			name:     "violates check constraint",
			err:      errors.New("violates check constraint"),
			expected: true,
		},
		{
			name:     "invalid video ID format",
			err:      errors.New("invalid video ID format"),
			expected: true,
		},
		{
			name:     "Invalid message format",
			err:      errors.New("Invalid message format"),
			expected: true,
		},
		{
			name:     "non-retryable error",
			err:      errors.New("non-retryable error"),
			expected: true,
		},
		{
			name:     "regular error",
			err:      errors.New("regular error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNonRetryableError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewMessageHandler(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := NewMessageHandler(mockUC)

	assert.NotNil(t, handler)
	assert.Equal(t, mockUC, handler.orchestrateUC)
}

func TestMessageHandler_HandleMessage_ProcessedMessage_Trim(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := NewMessageHandler(mockUC)

	processedMsg := VideoProcessedMessage{
		VideoID:    "123",
		Filename:   "test.mp4",
		BucketPath: "trim/test.mp4",
		Status:     "completed",
	}

	msgBytes, _ := json.Marshal(processedMsg)
	mockUC.On("HandleTrimCompleted", "123", "test.mp4").Return(nil)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_ProcessedMessage_Edit(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := NewMessageHandler(mockUC)

	processedMsg := VideoProcessedMessage{
		VideoID:    "123",
		Filename:   "test.mp4",
		BucketPath: "edit/test.mp4",
		Status:     "completed",
	}

	msgBytes, _ := json.Marshal(processedMsg)
	mockUC.On("HandleEditCompleted", "123", "test.mp4").Return(nil)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_ProcessedMessage_AudioRemoval(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := NewMessageHandler(mockUC)

	processedMsg := VideoProcessedMessage{
		VideoID:    "123",
		Filename:   "test.mp4",
		BucketPath: "audio-removal/test.mp4",
		Status:     "completed",
	}

	msgBytes, _ := json.Marshal(processedMsg)
	mockUC.On("HandleAudioRemovalCompleted", "123", "test.mp4").Return(nil)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_ProcessedMessage_Watermarking(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := NewMessageHandler(mockUC)

	processedMsg := VideoProcessedMessage{
		VideoID:    "123",
		Filename:   "test.mp4",
		BucketPath: "watermarking/test.mp4",
		Status:     "completed",
	}

	msgBytes, _ := json.Marshal(processedMsg)
	mockUC.On("HandleWatermarkingCompleted", "123", "test.mp4").Return(nil)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_ProcessedMessage_GossipOpenClose(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := NewMessageHandler(mockUC)

	processedMsg := VideoProcessedMessage{
		VideoID:    "123",
		Filename:   "test.mp4",
		BucketPath: "processed-videos/test.mp4",
		Status:     "completed",
	}

	msgBytes, _ := json.Marshal(processedMsg)
	mockUC.On("HandleGossipOpenCloseCompleted", "123", "test.mp4").Return(nil)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_ProcessedMessage_InvalidVideoID(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := NewMessageHandler(mockUC)

	processedMsg := VideoProcessedMessage{
		VideoID:    "123",
		Filename:   "test.mp4",
		BucketPath: "trim/test.mp4",
		Status:     "completed",
	}

	msgBytes, _ := json.Marshal(processedMsg)
	mockUC.On("HandleTrimCompleted", "123", "test.mp4").Return(errors.New("invalid video ID format"))

	err := handler.HandleMessage(msgBytes)

	assert.Error(t, err)
	assert.IsType(t, &NonRetryableError{}, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_InvalidJSON(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := NewMessageHandler(mockUC)

	invalidJSON := []byte("invalid json")

	err := handler.HandleMessage(invalidJSON)

	assert.Error(t, err)
	assert.IsType(t, &NonRetryableError{}, err)
	assert.Contains(t, err.Error(), "Invalid message format")
}

func TestMessageHandler_HandleMessage_MaxRetriesExceeded(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := NewMessageHandler(mockUC)

	msg := VideoMessage{
		VideoID:    "123",
		Filename:   "test.mp4",
		RetryCount: 5,
		MaxRetries: 3,
	}

	msgBytes, _ := json.Marshal(msg)

	err := handler.HandleMessage(msgBytes)

	assert.Error(t, err)
	assert.IsType(t, &NonRetryableError{}, err)
	assert.Contains(t, err.Error(), "Max retries exceeded")
}

func TestMessageHandler_HandleMessage_RetryDelayNotMet(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := NewMessageHandler(mockUC)

	msg := VideoMessage{
		VideoID:   "123",
		Filename:  "test.mp4",
		LastRetry: time.Now().Unix() - 60, // 1 minute ago
	}

	msgBytes, _ := json.Marshal(msg)
	mockUC.On("GetRetryDelayMinutes").Return(5) // 5 minutes required

	err := handler.HandleMessage(msgBytes)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retry delay not met")
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_Success(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := NewMessageHandler(mockUC)

	msg := VideoMessage{
		VideoID:  "123",
		Filename: "test.mp4",
	}

	msgBytes, _ := json.Marshal(msg)
	mockUC.On("Execute", "123").Return(nil)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_ExecuteError_InvalidVideoID(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := NewMessageHandler(mockUC)

	msg := VideoMessage{
		VideoID:  "123",
		Filename: "test.mp4",
	}

	msgBytes, _ := json.Marshal(msg)
	mockUC.On("Execute", "123").Return(errors.New("invalid video ID format"))

	err := handler.HandleMessage(msgBytes)

	assert.Error(t, err)
	assert.IsType(t, &NonRetryableError{}, err)
	mockUC.AssertExpectations(t)
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{"contains substring", "hello world", "world", true},
		{"does not contain", "hello world", "foo", false},
		{"empty substring", "hello", "", true},
		{"empty string", "", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVideoMessage_Structure(t *testing.T) {
	msg := VideoMessage{
		VideoID:     "123",
		Filename:    "test.mp4",
		RetryCount:  1,
		MaxRetries:  3,
		LastRetry:   1234567890,
	}

	assert.Equal(t, "123", msg.VideoID)
	assert.Equal(t, "test.mp4", msg.Filename)
	assert.Equal(t, 1, msg.RetryCount)
	assert.Equal(t, 3, msg.MaxRetries)
	assert.Equal(t, int64(1234567890), msg.LastRetry)
}

func TestVideoProcessedMessage_Structure(t *testing.T) {
	msg := VideoProcessedMessage{
		VideoID:     "123",
		Filename:    "test.mp4",
		BucketPath:  "trim/test.mp4",
		Status:      "completed",
		RetryCount:  1,
		MaxRetries:  3,
		LastRetry:   1234567890,
	}

	assert.Equal(t, "123", msg.VideoID)
	assert.Equal(t, "test.mp4", msg.Filename)
	assert.Equal(t, "trim/test.mp4", msg.BucketPath)
	assert.Equal(t, "completed", msg.Status)
	assert.Equal(t, 1, msg.RetryCount)
	assert.Equal(t, 3, msg.MaxRetries)
	assert.Equal(t, int64(1234567890), msg.LastRetry)
}