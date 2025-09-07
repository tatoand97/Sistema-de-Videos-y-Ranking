package adapters

import (
	"statesmachine/internal/application/usecases"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"strings"
	"shared/security"
)

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
		
		if contains(processedMsg.BucketPath, "trim") {
			return h.orchestrateUC.HandleTrimCompleted(processedMsg.VideoID, processedMsg.Filename)
		} else if contains(processedMsg.BucketPath, "edit") {
			return h.orchestrateUC.HandleEditCompleted(processedMsg.VideoID, processedMsg.Filename)
		} else if contains(processedMsg.BucketPath, "audio-removal") {
			return h.orchestrateUC.HandleAudioRemovalCompleted(processedMsg.VideoID, processedMsg.Filename)
		} else if contains(processedMsg.BucketPath, "watermarking") {
			return h.orchestrateUC.HandleWatermarkingCompleted(processedMsg.VideoID, processedMsg.Filename)
		} else if contains(processedMsg.BucketPath, "processed-videos") {
			return h.orchestrateUC.HandleGossipOpenCloseCompleted(processedMsg.VideoID, processedMsg.Filename)
		}
	}

	var msg VideoMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	logrus.Infof("StatesMachine received videoId: '%s'", security.SanitizeLogInput(msg.VideoID))
	return h.orchestrateUC.Execute(msg.VideoID)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

