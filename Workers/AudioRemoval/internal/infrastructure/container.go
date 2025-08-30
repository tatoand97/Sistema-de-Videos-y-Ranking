package infrastructure

import (
	"audioremoval/internal/adapters"
	"audioremoval/internal/application/services"
	"audioremoval/internal/application/usecases"
	"audioremoval/internal/domain"
)

type Container struct {
	Config              *Config
	VideoRepo           domain.VideoRepository
	StorageRepo         domain.StorageRepository
	ProcessingService   domain.VideoProcessingService
	NotificationService domain.NotificationService
	ProcessVideoUC      *usecases.ProcessVideoUseCase
}

func NewContainer(config *Config) (*Container, error) {
	storage, err := adapters.NewMinIOStorage(
		config.MinIOEndpoint,
		config.MinIOAccessKey,
		config.MinIOSecretKey,
	)
	if err != nil {
		return nil, err
	}

	videoRepo := adapters.NewVideoRepository()
	storageRepo := adapters.NewStorageRepository(storage)
	processingService := services.NewMP4VideoProcessingService()
	notificationService := services.NewLogNotificationService()

	processVideoUC := usecases.NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		config.RawBucket,
		config.ProcessedBucket,
	)

	return &Container{
		Config:              config,
		VideoRepo:           videoRepo,
		StorageRepo:         storageRepo,
		ProcessingService:   processingService,
		NotificationService: notificationService,
		ProcessVideoUC:      processVideoUC,
	}, nil
}