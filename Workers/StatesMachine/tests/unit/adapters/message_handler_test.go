package adapters

import (
	"encoding/json"
	"errors"
	"testing"
	"statesmachine/internal/adapters"

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
	nonRetryable := &adapters.NonRetryableError{
		OriginalError: originalErr,
		Message:       "test message",
	}

	expected := "non-retryable error: test message - original error"
	assert.Equal(t, expected, nonRetryable.Error())
}

func TestNonRetryableError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	nonRetryable := &adapters.NonRetryableError{
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
			err:      &adapters.NonRetryableError{OriginalError: errors.New("test"), Message: "test"},
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
			result := adapters.IsNonRetryableError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewMessageHandler(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	assert.NotNil(t, handler)
}

func TestMessageHandler_HandleMessage_ProcessedMessage_Trim(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	processedMsg := adapters.VideoProcessedMessage{
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

func TestMessageHandler_HandleMessage_InvalidJSON(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	invalidJSON := []byte("invalid json")

	err := handler.HandleMessage(invalidJSON)

	assert.Error(t, err)
	assert.IsType(t, &adapters.NonRetryableError{}, err)
	assert.Contains(t, err.Error(), "Invalid message format")
}

func TestMessageHandler_HandleMessage_Success(t *testing.T) {
	mockUC := &MockOrchestrateUseCase{}
	handler := adapters.NewMessageHandler(mockUC)

	msg := adapters.VideoMessage{
		VideoID:  "123",
		Filename: "test.mp4",
	}

	msgBytes, _ := json.Marshal(msg)
	mockUC.On("Execute", "123").Return(nil)

	err := handler.HandleMessage(msgBytes)

	assert.NoError(t, err)
	mockUC.AssertExpectations(t)
}