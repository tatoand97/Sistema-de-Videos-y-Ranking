package adapters

import (
	"gossipopenclose/internal/application/usecases"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type VideoMessage struct {
	Filename string `json:"filename"`
}

type MessageHandler struct {
	editVideoUC *usecases.OpenCloseUseCase
}

func NewMessageHandler(uc *usecases.OpenCloseUseCase) *MessageHandler {
	return &MessageHandler{editVideoUC: uc}
}

func (h *MessageHandler) HandleMessage(body []byte) error {
	var msg VideoMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	logrus.Infof("Recibido filename: '%s'", msg.Filename)
	return h.editVideoUC.Execute(msg.Filename)
}
