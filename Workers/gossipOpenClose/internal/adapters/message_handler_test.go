package adapters

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOpenCloseUseCase struct {
	mock.Mock
}

func (m *MockOpenCloseUseCase) Execute(videoID, filename string) error {
	args := m.Called(videoID, filename)
	return args.Error(0)
}

func TestNewMessageHandler(t *testing.T) {
	useCase := &MockOpenCloseUseCase{}
	handler := NewMessageHandler(useCase)
	
	assert.NotNil(t, handler)
	assert.Equal(t, useCase, handler.editVideoUC)
}

func TestMessageHandler_HandleMessage_Success(t *testing.T) {
	useCase := &MockOpenCloseUseCase{}
	handler := NewMessageHandler(useCase)
	
	msg := VideoMessage{
		VideoID:  "video-123",
		Filename: "test.mp4",
	}
	
	body, _ := json.Marshal(msg)
	useCase.On("Execute", "video-123", "test.mp4").Return(nil)
	
	err := handler.HandleMessage(body)
	
	assert.NoError(t, err)
	useCase.AssertExpectations(t)
}

func TestMessageHandler_HandleMessage_InvalidJSON(t *testing.T) {
	useCase := &MockOpenCloseUseCase{}
	handler := NewMessageHandler(useCase)
	
	invalidJSON := []byte(`{"invalid": json}`)
	
	err := handler.HandleMessage(invalidJSON)
	
	assert.Error(t, err)
	useCase.AssertNotCalled(t, "Execute")
}

func TestMessageHandler_HandleMessage_ExecuteError(t *testing.T) {
	useCase := &MockOpenCloseUseCase{}
	handler := NewMessageHandler(useCase)
	
	msg := VideoMessage{
		VideoID:  "video-123",
		Filename: "test.mp4",
	}
	
	body, _ := json.Marshal(msg)
	expectedError := errors.New("execution failed")
	useCase.On("Execute", "video-123", "test.mp4").Return(expectedError)
	
	err := handler.HandleMessage(body)
	
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	useCase.AssertExpectations(t)
}

func TestVideoMessage_Structure(t *testing.T) {
	msg := VideoMessage{
		VideoID:  "test-id",
		Filename: "test-file.mp4",
	}
	
	assert.Equal(t, "test-id", msg.VideoID)
	assert.Equal(t, "test-file.mp4", msg.Filename)
}

func TestMessageHandler_HandleMessage_EmptyFields(t *testing.T) {
	useCase := &MockOpenCloseUseCase{}
	handler := NewMessageHandler(useCase)
	
	tests := []struct {
		name     string
		msg      VideoMessage
		shouldCall bool
	}{
		{
			name: "empty video ID",
			msg: VideoMessage{
				VideoID:  "",
				Filename: "test.mp4",
			},
			shouldCall: true,
		},
		{
			name: "empty filename",
			msg: VideoMessage{
				VideoID:  "video-123",
				Filename: "",
			},
			shouldCall: true,
		},
		{
			name: "both empty",
			msg: VideoMessage{
				VideoID:  "",
				Filename: "",
			},
			shouldCall: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.msg)
			
			if tt.shouldCall {
				useCase.On("Execute", tt.msg.VideoID, tt.msg.Filename).Return(nil).Once()
			}
			
			err := handler.HandleMessage(body)
			
			if tt.shouldCall {
				assert.NoError(t, err)
			}
		})
	}
	
	useCase.AssertExpectations(t)
}