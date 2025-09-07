package adapters

import (
	"statesmachine/internal/application/usecases"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"shared/security"
)

// NonRetryableError represents errors that should not be retried
type NonRetryableError struct {
	OriginalError error
	Message       string
}

func (e *NonRetryableError) Error() string {
	return fmt.Sprintf("non-retryable error: %s - %v", e.Message, e.OriginalError)
}

func (e *NonRetryableError) Unwrap() error {
	return e.OriginalError
}

func IsNonRetryableError(err error) bool {
	var nonRetryable *NonRetryableError
	if errors.As(err, &nonRetryable) {
		return true
	}
	
	// Check for database constraint errors
	errorMsg := err.Error()
	return strings.Contains(errorMsg, "violates check constraint") ||
		   strings.Contains(errorMsg, "invalid video ID format") ||
		   strings.Contains(errorMsg, "Invalid message format") ||
		   strings.Contains(errorMsg, "non-retryable error")
}

type VideoMessage struct {
	VideoID  string `json:"videoId"`
	Filename string `json:"filename"`
}

type VideoProcessedMessage struct {
	VideoID    string `json:"video_id"`
	Filename   string `json:"filename"`
	BucketPath string `json:"bucket_path"`
	Status     string `json:"status"`
}

type MessageHandler struct {
	orchestrateUC *usecases.OrchestrateVideoUseCase
}

func NewMessageHandler(uc *usecases.OrchestrateVideoUseCase) *MessageHandler {
	return &MessageHandler{orchestrateUC: uc}
}

func (h *MessageHandler) HandleMessage(body []byte) error {
	var processedMsg VideoProcessedMessage
	if err := json.Unmarshal(body, &processedMsg); err == nil && processedMsg.VideoID != "" {
		logrus.Infof("StatesMachine received processed video: %s from %s", security.SanitizeLogInput(processedMsg.Filename), security.SanitizeLogInput(processedMsg.BucketPath))
		
		var handlerErr error
		if contains(processedMsg.BucketPath, "trim") {
			handlerErr = h.orchestrateUC.HandleTrimCompleted(processedMsg.VideoID, processedMsg.Filename)
		} else if contains(processedMsg.BucketPath, "edit") {
			handlerErr = h.orchestrateUC.HandleEditCompleted(processedMsg.VideoID, processedMsg.Filename)
		} else if contains(processedMsg.BucketPath, "audio-removal") {
			handlerErr = h.orchestrateUC.HandleAudioRemovalCompleted(processedMsg.VideoID, processedMsg.Filename)
		} else if contains(processedMsg.BucketPath, "watermarking") {
			handlerErr = h.orchestrateUC.HandleWatermarkingCompleted(processedMsg.VideoID, processedMsg.Filename)
		} else if contains(processedMsg.BucketPath, "processed-videos") {
			handlerErr = h.orchestrateUC.HandleGossipOpenCloseCompleted(processedMsg.VideoID, processedMsg.Filename)
		}
		
		if handlerErr != nil && strings.Contains(handlerErr.Error(), "invalid video ID format") {
			return &NonRetryableError{
				OriginalError: handlerErr,
				Message:       "Invalid video ID in processed message",
			}
		}
		return handlerErr
	}

	var msg VideoMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return &NonRetryableError{
			OriginalError: err,
			Message:       "Invalid message format",
		}
	}

	logrus.Infof("StatesMachine received videoId: '%s'", security.SanitizeLogInput(msg.VideoID))
	execErr := h.orchestrateUC.Execute(msg.VideoID)
	if execErr != nil && strings.Contains(execErr.Error(), "invalid video ID format") {
		return &NonRetryableError{
			OriginalError: execErr,
			Message:       "Invalid video ID in orchestration message",
		}
	}
	return execErr
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

