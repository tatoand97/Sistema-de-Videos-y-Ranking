package adapters

import (
	"editvideo/internal/application/usecases"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

type VideoMessage struct {
	Filename string `json:"filename"`
}

type MessageHandler struct {
	editVideoUC *usecases.EditVideoUseCase
}

func NewMessageHandler(uc *usecases.EditVideoUseCase) *MessageHandler {
	return &MessageHandler{editVideoUC: uc}
}

func (h *MessageHandler) HandleMessage(body []byte) error {
	var msg VideoMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	logrus.Infof("Received filename: '%s'", sanitizeLogInput(msg.Filename))
	return h.editVideoUC.Execute(msg.Filename)
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
