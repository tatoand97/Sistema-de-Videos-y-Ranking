package adapters

import (
	"audioremoval/internal/application/usecases"
	"encoding/json"
	"github.com/sirupsen/logrus"
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
	logrus.Infof("Received message: %s", string(body))
	
	var msg VideoMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	logrus.Infof("Parsed filename: '%s'", msg.Filename)
	logrus.Infof("Processing video: %s", msg.Filename)
	
	if err := h.processVideoUC.Execute(msg.Filename); err != nil {
		logrus.Errorf("Error processing video %s: %v", msg.Filename, err)
		return err
	}

	logrus.Infof("Video processed successfully: %s", msg.Filename)
	return nil
}