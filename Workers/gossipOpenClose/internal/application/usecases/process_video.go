
package usecases

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gossipopenclose/internal/domain"
	"gossipopenclose/internal/application/services"
	"github.com/sirupsen/logrus"
)

type OpenCloseUseCase struct {
	videoRepo            domain.VideoRepository
	storageRepo          domain.StorageRepository
	processingService    *services.OpenCloseVideoProcessingService
	notificationService  domain.NotificationService
	rawBucket            string
	processedBucket      string
	introSeconds         float64
	outroSeconds         float64
	targetW              int
	targetH              int
	fps                  int
	logoPath             string
}

func NewOpenCloseUseCase(
	videoRepo domain.VideoRepository,
	storageRepo domain.StorageRepository,
	processingService *services.OpenCloseVideoProcessingService,
	notificationService domain.NotificationService,
	rawBucket, processedBucket, logoPath string,
	introSeconds, outroSeconds float64,
	targetW, targetH, fps int,
) *OpenCloseUseCase {
	return &OpenCloseUseCase{
		videoRepo: videoRepo,
		storageRepo: storageRepo,
		processingService: processingService,
		notificationService: notificationService,
		rawBucket: rawBucket,
		processedBucket: processedBucket,
		introSeconds: introSeconds,
		outroSeconds: outroSeconds,
		targetW: targetW,
		targetH: targetH,
		fps: fps,
		logoPath: logoPath,
	}
}

func envFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil { return f }
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil { return i }
	}
	return def
}

func (uc *OpenCloseUseCase) Execute(filename string) error {
	video, err := uc.videoRepo.FindByFilename(filename)
	if err != nil { return fmt.Errorf("find video: %w", err) }

	if err := uc.videoRepo.UpdateStatus(video.ID, domain.StatusProcessing); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	data, err := uc.storageRepo.Download(uc.rawBucket, filename)
	if err != nil {
		_ = uc.videoRepo.UpdateStatus(video.ID, domain.StatusFailed)
		return fmt.Errorf("download: %w", err)
	}

	processedData, err := uc.processingService.Process(data, uc.logoPath, uc.introSeconds, uc.outroSeconds, uc.targetW, uc.targetH, uc.fps)
	if err != nil {
		_ = uc.videoRepo.UpdateStatus(video.ID, domain.StatusFailed)
		return fmt.Errorf("process: %w", err)
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
		logrus.Errorf("Failed to notify state machine: %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"video_id": video.ID,
		"filename": filename,
		"bucket_from": uc.rawBucket,
		"bucket_to": uc.processedBucket,
		"timestamp": time.Now().UTC(),
	}).Info("GossipOpenClose processing completed successfully")

	return nil
}
