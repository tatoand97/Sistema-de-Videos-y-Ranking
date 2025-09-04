package adapters

import (
	"statesmachine/internal/application/usecases"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type VideoMessage struct {
	Filename string `json:"filename"`
}

type MessageHandler struct {
	orchestrateUC *usecases.OrchestrateVideoUseCase
}

func NewMessageHandler(uc *usecases.OrchestrateVideoUseCase) *MessageHandler {
	return &MessageHandler{orchestrateUC: uc}
}

func (h *MessageHandler) HandleMessage(body []byte) error {
	var msg VideoMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	logrus.Infof("StatesMachine received filename: '%s'", msg.Filename)
	return h.orchestrateUC.Execute(msg.Filename)
}