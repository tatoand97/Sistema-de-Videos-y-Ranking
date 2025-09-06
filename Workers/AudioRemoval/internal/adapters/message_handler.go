package adapters

import (
	"audioremoval/internal/application/usecases"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

type VideoMessage struct {
	Filename string `json:"filename"`
}

type MessageHandler struct {
	processVideoUC *usecases.ProcessVideoUseCase
}

func NewMessageHandler(processVideoUC *usecases.ProcessVideoUseCase) *MessageHandler {
	return &MessageHandler{processVideoUC: processVideoUC}
}

func (h *MessageHandler) HandleMessage(body []byte) error {
	logrus.Infof("Received message: %s", sanitizeLogInput(string(body)))
	
	var msg VideoMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	logrus.Infof("Parsed filename: '%s'", sanitizeLogInput(msg.Filename))
	logrus.Infof("Processing video: %s", sanitizeLogInput(msg.Filename))
	
	if err := h.processVideoUC.Execute(msg.Filename); err != nil {
		logrus.Errorf("Error processing video %s: %v", sanitizeLogInput(msg.Filename), err)
		return err
	}

	logrus.Infof("Video processed successfully: %s", sanitizeLogInput(msg.Filename))
	return nil
}

// sanitizeLogInput removes potentially dangerous characters from log input
func sanitizeLogInput(input string) string {
	// Remove newlines, carriage returns, and control characters
	re := regexp.MustCompile(`[\r\n\t\x00-\x1f\x7f-\x9f]`)
	sanitized := re.ReplaceAllString(input, "")
	// Limit length to prevent log flooding
	if len(sanitized) > 100 {
		sanitized = sanitized[:100] + "..."
	}
	return strings.TrimSpace(sanitized)
}