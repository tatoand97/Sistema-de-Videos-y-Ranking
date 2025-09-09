package adapters

import (
	"watermarking/internal/application/usecases"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type VideoMessage struct {
	VideoID  string `json:"video_id"`
	Filename string `json:"filename"`
}

type MessageHandler struct {
	editVideoUC *usecases.WatermarkingUseCase
}

func NewMessageHandler(uc *usecases.WatermarkingUseCase) *MessageHandler {
	return &MessageHandler{editVideoUC: uc}
}

func (h *MessageHandler) HandleMessage(body []byte) error {
	var msg VideoMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	logrus.Infof("Recibido video_id: '%s', filename: '%s'", msg.VideoID, msg.Filename)
	return h.editVideoUC.Execute(msg.VideoID, msg.Filename)
}
