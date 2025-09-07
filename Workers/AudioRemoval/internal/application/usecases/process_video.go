package usecases

import (
	"audioremoval/internal/domain"
	"fmt"
	"time"
	"github.com/sirupsen/logrus"
)

type ProcessVideoUseCase struct {
	videoRepo       domain.VideoRepository
	storageRepo     domain.StorageRepository
	processingService domain.VideoProcessingService
	notificationService domain.NotificationService
	rawBucket       string
	processedBucket string
}

func NewProcessVideoUseCase(
	videoRepo domain.VideoRepository,
	storageRepo domain.StorageRepository,
	processingService domain.VideoProcessingService,
	notificationService domain.NotificationService,
	rawBucket, processedBucket string,
) *ProcessVideoUseCase {
	return &ProcessVideoUseCase{
		videoRepo:           videoRepo,
		storageRepo:         storageRepo,
		processingService:   processingService,
		notificationService: notificationService,
		rawBucket:           rawBucket,
		processedBucket:     processedBucket,
	}
}

func (uc *ProcessVideoUseCase) Execute(videoID, filename string) error {
	video, err := uc.videoRepo.FindByFilename(filename)
	if err != nil {
		return fmt.Errorf("video not found: %w", err)
	}

	if err := uc.videoRepo.UpdateStatus(video.ID, domain.StatusProcessing); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	inputData, err := uc.storageRepo.Download(uc.rawBucket, filename)
	if err != nil {
		uc.videoRepo.UpdateStatus(video.ID, domain.StatusFailed)
		return fmt.Errorf("failed to download video: %w", err)
	}

	processedData, err := uc.processingService.RemoveAudio(inputData)
	if err != nil {
		uc.videoRepo.UpdateStatus(video.ID, domain.StatusFailed)
		return fmt.Errorf("failed to process video: %w", err)
	}

	if err := uc.storageRepo.Upload(uc.processedBucket, filename, processedData); err != nil {
		uc.videoRepo.UpdateStatus(video.ID, domain.StatusFailed)
		return fmt.Errorf("failed to upload processed video: %w", err)
	}

	if err := uc.videoRepo.UpdateStatus(video.ID, domain.StatusCompleted); err != nil {
		return fmt.Errorf("failed to update final status: %w", err)
	}

	bucketPath := fmt.Sprintf("%s/%s", uc.processedBucket, filename)
	if err := uc.notificationService.NotifyVideoProcessed(videoID, filename, bucketPath); err != nil {
		logrus.Errorf("Failed to notify state machine: %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"video_id": video.ID,
		"filename": filename,
		"bucket_from": uc.rawBucket,
		"bucket_to": uc.processedBucket,
		"timestamp": time.Now().UTC(),
	}).Info("AudioRemoval processing completed successfully")

	return nil
}