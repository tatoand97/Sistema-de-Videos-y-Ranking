package adapters

import (
	"trimvideo/internal/application/usecases"
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

func NewMessageHandler(uc *usecases.ProcessVideoUseCase) *MessageHandler {
	return &MessageHandler{processVideoUC: uc}
}

func (h *MessageHandler) HandleMessage(body []byte) error {
	var msg VideoMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	logrus.Infof("Received video_id: '%s', filename: '%s'", security.SanitizeLogInput(msg.VideoID), security.SanitizeLogInput(msg.Filename))
	return h.processVideoUC.Execute(msg.VideoID, msg.Filename)
}




