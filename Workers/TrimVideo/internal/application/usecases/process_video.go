package usecases

import (
    "fmt"
    "time"
    "trimvideo/internal/domain"
    "github.com/sirupsen/logrus"
)

type ProcessVideoUseCase struct {
	videoRepo            domain.VideoRepository
	storageRepo          domain.StorageRepository
	processingService    domain.VideoProcessingService
	notificationService  domain.NotificationService
	rawBucket            string
	processedBucket      string
	maxSeconds           int
}

func NewProcessVideoUseCase(
	videoRepo domain.VideoRepository,
	storageRepo domain.StorageRepository,
	processingService domain.VideoProcessingService,
	notificationService domain.NotificationService,
	rawBucket, processedBucket string,
	maxSeconds int,
) *ProcessVideoUseCase {
	return &ProcessVideoUseCase{
		videoRepo: videoRepo,
		storageRepo: storageRepo,
		processingService: processingService,
		notificationService: notificationService,
		rawBucket: rawBucket,
		processedBucket: processedBucket,
		maxSeconds: maxSeconds,
	}
}

func (uc *ProcessVideoUseCase) Execute(videoID, filename string) error {
	video, err := uc.videoRepo.FindByFilename(filename)
	if err != nil { return fmt.Errorf("find video: %w", err) }
	if err := uc.videoRepo.UpdateStatus(video.ID, domain.StatusProcessing); err != nil { return err }

	data, err := uc.storageRepo.Download(uc.rawBucket, filename)
	if err != nil {
		_ = uc.videoRepo.UpdateStatus(video.ID, domain.StatusFailed)
		return fmt.Errorf("download: %w", err)
	}

	processedData, err := uc.processingService.TrimToMaxSeconds(data, uc.maxSeconds)
	if err != nil {
		_ = uc.videoRepo.UpdateStatus(video.ID, domain.StatusFailed)
		return fmt.Errorf("processing: %w", err)
	}

	if err := uc.storageRepo.Upload(uc.processedBucket, filename, processedData); err != nil {
		_ = uc.videoRepo.UpdateStatus(video.ID, domain.StatusFailed)
		return fmt.Errorf("upload: %w", err)
	}

	if err := uc.videoRepo.UpdateStatus(video.ID, domain.StatusCompleted); err != nil {
		return fmt.Errorf("update final status: %w", err)
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
		"max_seconds": uc.maxSeconds,
		"timestamp": time.Now().UTC(),
	}).Info("TrimVideo processing completed successfully")

	return nil
}
