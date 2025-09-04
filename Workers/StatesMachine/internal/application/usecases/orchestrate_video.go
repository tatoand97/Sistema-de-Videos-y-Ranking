package usecases

import (
	"encoding/json"
	"fmt"
	"time"
	"statesmachine/internal/domain"
	"github.com/sirupsen/logrus"
)

type VideoMessage struct {
	Filename string `json:"filename"`
}

type OrchestrateVideoUseCase struct {
	videoRepo domain.VideoRepository
	publisher domain.MessagePublisher
}

func NewOrchestrateVideoUseCase(
	videoRepo domain.VideoRepository,
	publisher domain.MessagePublisher,
) *OrchestrateVideoUseCase {
	return &OrchestrateVideoUseCase{
		videoRepo: videoRepo,
		publisher: publisher,
	}
}

func (uc *OrchestrateVideoUseCase) Execute(filename string) error {
	logrus.WithFields(logrus.Fields{
		"filename":  filename,
		"timestamp": time.Now().UTC(),
		"stage":     "orchestration_start",
	}).Info("StatesMachine: Starting video processing orchestration")

	video, err := uc.videoRepo.FindByFilename(filename)
	if err != nil {
		return fmt.Errorf("find video: %w", err)
	}

	if err := uc.videoRepo.UpdateStatus(video.ID, domain.StatusProcessing); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	message := VideoMessage{Filename: filename}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if err := uc.publisher.PublishMessage("trim_video_queue", messageBytes); err != nil {
		_ = uc.videoRepo.UpdateStatus(video.ID, domain.StatusFailed)
		return fmt.Errorf("publish to trim_video_queue: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"filename":    filename,
		"video_id":    video.ID,
		"next_queue":  "trim_video_queue",
		"timestamp":   time.Now().UTC(),
	}).Info("StatesMachine: Message published to TrimVideo queue")

	return nil
}