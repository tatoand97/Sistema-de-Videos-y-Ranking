package usecases

import (
    "fmt"
    "editvideo/internal/domain"
)

type EditVideoUseCase struct {
	videoRepo            domain.VideoRepository
	storageRepo          domain.StorageRepository
	processingService    domain.VideoProcessingService
	notificationService  domain.NotificationService
	rawBucket            string
	processedBucket      string
	maxSeconds           int
}

func NewEditVideoUseCase(
	videoRepo domain.VideoRepository,
	storageRepo domain.StorageRepository,
	processingService domain.VideoProcessingService,
	notificationService domain.NotificationService,
	rawBucket, processedBucket string,
	maxSeconds int,
) *EditVideoUseCase {
	return &EditVideoUseCase{
		videoRepo: videoRepo,
		storageRepo: storageRepo,
		processingService: processingService,
		notificationService: notificationService,
		rawBucket: rawBucket,
		processedBucket: processedBucket,
		maxSeconds: maxSeconds,
	}
}

func (uc *EditVideoUseCase) Execute(filename string) error {
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
	if err := uc.notificationService.NotifyVideoProcessed(video.ID, filename, bucketPath); err != nil {
		return fmt.Errorf("notify state machine: %w", err)
	}

	return nil
}
