package usecases

import (
	"encoding/json"
	"fmt"
	"time"
	"statesmachine/internal/domain"
	"github.com/sirupsen/logrus"
)

type WorkerMessage struct {
	VideoID  string `json:"video_id"`
	Filename string `json:"filename"`
}

type OrchestrateVideoUseCase struct {
	videoRepo         domain.VideoRepository
	publisher         domain.MessagePublisher
	editVideoQueue    string
	audioRemovalQueue string
	watermarkingQueue string
}

func NewOrchestrateVideoUseCase(
	videoRepo domain.VideoRepository,
	publisher domain.MessagePublisher,
	editVideoQueue, audioRemovalQueue, watermarkingQueue string,
) *OrchestrateVideoUseCase {
	return &OrchestrateVideoUseCase{
		videoRepo:         videoRepo,
		publisher:         publisher,
		editVideoQueue:    editVideoQueue,
		audioRemovalQueue: audioRemovalQueue,
		watermarkingQueue: watermarkingQueue,
	}
}

func (uc *OrchestrateVideoUseCase) Execute(videoID string) error {
	logrus.WithFields(logrus.Fields{
		"video_id":  videoID,
		"timestamp": time.Now().UTC(),
		"stage":     "orchestration_start",
	}).Info("StatesMachine: Starting video processing orchestration")

	// Convert string ID to uint
	var id uint
	if _, err := fmt.Sscanf(videoID, "%d", &id); err != nil {
		return fmt.Errorf("invalid video ID format: %w", err)
	}

	video, err := uc.videoRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("find video: %w", err)
	}

	message := WorkerMessage{
		VideoID:  videoID,
		Filename: video.OriginalFile,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if err := uc.publisher.PublishMessage("trim_video_queue", messageBytes); err != nil {
		return fmt.Errorf("publish to trim_video_queue: %w", err)
	}

	if err := uc.videoRepo.UpdateStatus(video.ID, domain.StatusTrimming); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"video_id":    videoID,
		"filename":    video.OriginalFile,
		"next_queue":  "trim_video_queue",
		"timestamp":   time.Now().UTC(),
	}).Info("StatesMachine: Message published to TrimVideo queue")

	return nil
}

func (uc *OrchestrateVideoUseCase) HandleTrimCompleted(videoID, filename string) error {
	logrus.WithFields(logrus.Fields{
		"video_id":  videoID,
		"filename":  filename,
		"timestamp": time.Now().UTC(),
		"stage":     "trim_completed",
	}).Info("StatesMachine: TrimVideo completed, sending to EditVideo")

	message := WorkerMessage{
		VideoID:  videoID,
		Filename: filename,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if err := uc.publisher.PublishMessage(uc.editVideoQueue, messageBytes); err != nil {
		return fmt.Errorf("publish to edit_video_queue: %w", err)
	}

	// Update status to ADJUSTING_RESOLUTION
	var id uint
	if _, err := fmt.Sscanf(videoID, "%d", &id); err != nil {
		return fmt.Errorf("invalid video ID format: %w", err)
	}
	if err := uc.videoRepo.UpdateStatus(id, domain.StatusAdjustingRes); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"video_id":   videoID,
		"filename":   filename,
		"next_queue": uc.editVideoQueue,
		"timestamp":  time.Now().UTC(),
	}).Info("StatesMachine: Message published to EditVideo queue")

	return nil
}

func (uc *OrchestrateVideoUseCase) HandleEditCompleted(videoID, filename string) error {
	logrus.WithFields(logrus.Fields{
		"video_id":  videoID,
		"filename":  filename,
		"timestamp": time.Now().UTC(),
		"stage":     "edit_completed",
		"result":    "success",
	}).Info("StatesMachine: EditVideo completed successfully, sending to AudioRemoval")

	message := WorkerMessage{
		VideoID:  videoID,
		Filename: filename,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if err := uc.publisher.PublishMessage(uc.audioRemovalQueue, messageBytes); err != nil {
		return fmt.Errorf("publish to audio_removal_queue: %w", err)
	}

	// Update status to REMOVING_AUDIO
	var id uint
	if _, err := fmt.Sscanf(videoID, "%d", &id); err != nil {
		return fmt.Errorf("invalid video ID format: %w", err)
	}
	if err := uc.videoRepo.UpdateStatus(id, domain.StatusRemovingAudio); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"video_id":   videoID,
		"filename":   filename,
		"next_queue": uc.audioRemovalQueue,
		"timestamp":  time.Now().UTC(),
	}).Info("StatesMachine: Message published to AudioRemoval queue")

	return nil
}

func (uc *OrchestrateVideoUseCase) HandleAudioRemovalCompleted(videoID, filename string) error {
	logrus.WithFields(logrus.Fields{
		"video_id":  videoID,
		"filename":  filename,
		"timestamp": time.Now().UTC(),
		"stage":     "audio_removal_completed",
		"result":    "success",
	}).Info("StatesMachine: AudioRemoval completed successfully, sending to Watermarking")

	message := WorkerMessage{
		VideoID:  videoID,
		Filename: filename,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if err := uc.publisher.PublishMessage(uc.watermarkingQueue, messageBytes); err != nil {
		return fmt.Errorf("publish to watermarking_queue: %w", err)
	}

	// Update status to ADDING_WATERMARK
	var id uint
	if _, err := fmt.Sscanf(videoID, "%d", &id); err != nil {
		return fmt.Errorf("invalid video ID format: %w", err)
	}
	if err := uc.videoRepo.UpdateStatus(id, domain.StatusAddingWatermark); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"video_id":   videoID,
		"filename":   filename,
		"next_queue": uc.watermarkingQueue,
		"timestamp":  time.Now().UTC(),
	}).Info("StatesMachine: Message published to Watermarking queue")

	return nil
}

func (uc *OrchestrateVideoUseCase) HandleWatermarkingCompleted(videoID, filename string) error {
	logrus.WithFields(logrus.Fields{
		"video_id":  videoID,
		"filename":  filename,
		"timestamp": time.Now().UTC(),
		"stage":     "watermarking_completed",
		"result":    "success",
	}).Info("StatesMachine: Watermarking completed successfully, sending to GossipOpenClose")

	message := WorkerMessage{
		VideoID:  videoID,
		Filename: filename,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if err := uc.publisher.PublishMessage("gossip_open_close_queue", messageBytes); err != nil {
		return fmt.Errorf("publish to gossip_open_close_queue: %w", err)
	}

	// Update status to ADDING_INTRO_OUTRO
	var id uint
	if _, err := fmt.Sscanf(videoID, "%d", &id); err != nil {
		return fmt.Errorf("invalid video ID format: %w", err)
	}
	if err := uc.videoRepo.UpdateStatus(id, domain.StatusAddingIntroOutro); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"video_id":   videoID,
		"filename":   filename,
		"next_queue": "gossip_open_close_queue",
		"timestamp":  time.Now().UTC(),
	}).Info("StatesMachine: Message published to GossipOpenClose queue")

	return nil
}

func (uc *OrchestrateVideoUseCase) HandleGossipOpenCloseCompleted(videoID, filename string) error {
	// Update status to PROCESSED
	var id uint
	if _, err := fmt.Sscanf(videoID, "%d", &id); err != nil {
		return fmt.Errorf("invalid video ID format: %w", err)
	}
	if err := uc.videoRepo.UpdateStatus(id, domain.StatusProcessed); err != nil {
		return fmt.Errorf("update final status: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"video_id":  videoID,
		"filename":  filename,
		"timestamp": time.Now().UTC(),
		"stage":     "gossip_open_close_completed",
		"result":    "success",
		"pipeline":  "finished",
	}).Info("StatesMachine: GossipOpenClose completed successfully, entire video processing pipeline finished")

	return nil
}