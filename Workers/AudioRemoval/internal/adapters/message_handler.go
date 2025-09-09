package adapters

import (
	"audioremoval/internal/application/usecases"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"shared/security"
)

type VideoMessage struct {
	VideoID  string `json:"video_id"`
	Filename string `json:"filename"`
}

type MessageHandler struct {
	processVideoUC *usecases.ProcessVideoUseCase
}

func NewMessageHandler(processVideoUC *usecases.ProcessVideoUseCase) *MessageHandler {
	return &MessageHandler{processVideoUC: processVideoUC}
}

func (h *MessageHandler) HandleMessage(body []byte) error {
	logrus.Infof("Received message: %s", security.SanitizeLogInput(string(body)))
	
	var msg VideoMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	logrus.Infof("Processing video_id: %s, filename: %s", security.SanitizeLogInput(msg.VideoID), security.SanitizeLogInput(msg.Filename))
	
	if err := h.processVideoUC.Execute(msg.VideoID, msg.Filename); err != nil {
		logrus.Errorf("Error processing video %s: %v", security.SanitizeLogInput(msg.Filename), err)
		return err
	}

	logrus.Infof("Video processed successfully: %s", security.SanitizeLogInput(msg.Filename))
	return nil
}

