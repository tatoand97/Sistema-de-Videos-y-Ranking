package usecases

import (
    "fmt"
    "mime"
    "time"
    "editvideo/internal/domain"
)

type EditVideoUseCase struct {
	videoRepo            domain.VideoRepository
	storageRepo          domain.StorageRepository
	processingService    domain.VideoProcessingService
	rawBucket            string
	processedBucket      string
	maxSeconds           int
}

func NewEditVideoUseCase(
	videoRepo domain.VideoRepository,
	storageRepo domain.StorageRepository,
	processingService domain.VideoProcessingService,
	rawBucket, processedBucket string,
	maxSeconds int,
) *EditVideoUseCase {
	return &EditVideoUseCase{
		videoRepo: videoRepo,
		storageRepo: storageRepo,
		processingService: processingService,
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
	_ = mime.TypeByExtension(".mp4") // placeholder like AudioRemoval; could be used for content-type

	_ = time.Now()
	return nil
}
