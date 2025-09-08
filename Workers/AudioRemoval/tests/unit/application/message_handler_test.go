package application_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/streadway/amqp"
)

type MessageHandler struct {
	usecase UseCaseInterface
	logger  Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type UseCaseInterface interface {
	ProcessVideo(ctx context.Context, videoID, inputURL string) (string, error)
}

type AudioRemovalMessage struct {
	VideoID   string `json:"video_id"`
	InputURL  string `json:"input_url"`
	OutputURL string `json:"output_url"`
}

func NewMessageHandler(usecase UseCaseInterface, logger Logger) *MessageHandler {
	return &MessageHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *MessageHandler) HandleMessage(ctx context.Context, delivery amqp.Delivery) error {
	var msg AudioRemovalMessage
	if err := json.Unmarshal(delivery.Body, &msg); err != nil {
		h.logger.Error("failed to unmarshal message: " + err.Error())
		return err
	}

	h.logger.Info("processing video: " + msg.VideoID)
	
	outputURL, err := h.usecase.ProcessVideo(ctx, msg.VideoID, msg.InputURL)
	if err != nil {
		h.logger.Error("failed to process video: " + err.Error())
		return err
	}

	h.logger.Info("video processed successfully: " + outputURL)
	return nil
}

type MockLogger struct {
	logs []string
}

func (m *MockLogger) Info(msg string) {
	m.logs = append(m.logs, "INFO: "+msg)
}

func (m *MockLogger) Error(msg string) {
	m.logs = append(m.logs, "ERROR: "+msg)
}

type MockUseCase struct {
	result string
	err    error
}

func (m *MockUseCase) ProcessVideo(ctx context.Context, videoID, inputURL string) (string, error) {
	return m.result, m.err
}

func TestMessageHandler_HandleMessage(t *testing.T) {
	tests := []struct {
		name     string
		message  AudioRemovalMessage
		mockUsecase *MockUseCase
		wantErr  bool
	}{
		{
			name: "successful processing",
			message: AudioRemovalMessage{
				VideoID:  "video-123",
				InputURL: "https://example.com/input.mp4",
			},
			mockUsecase: &MockUseCase{
				result: "https://example.com/output.mp4",
			},
			wantErr: false,
		},
		{
			name: "processing fails",
			message: AudioRemovalMessage{
				VideoID:  "video-123",
				InputURL: "https://example.com/input.mp4",
			},
			mockUsecase: &MockUseCase{
				err: errors.New("processing failed"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			handler := NewMessageHandler(tt.mockUsecase, mockLogger)

			msgBytes, _ := json.Marshal(tt.message)
			delivery := amqp.Delivery{Body: msgBytes}

			err := handler.HandleMessage(context.Background(), delivery)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMessageHandler_HandleMessage_InvalidJSON(t *testing.T) {
	mockUsecase := &MockUseCase{}
	mockLogger := &MockLogger{}
	handler := NewMessageHandler(mockUsecase, mockLogger)

	delivery := amqp.Delivery{Body: []byte("invalid json")}
	err := handler.HandleMessage(context.Background(), delivery)

	assert.Error(t, err)
	assert.Contains(t, mockLogger.logs[0], "ERROR:")
}