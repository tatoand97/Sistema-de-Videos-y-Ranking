package infrastructure

import (
	"trimvideo/internal/adapters"
	"trimvideo/internal/application/services"
	"trimvideo/internal/application/usecases"

	"shared/messaging"
	sharedstorage "shared/storage"
)

type Container struct {
	Config         *Config
	Consumer       *messaging.SQSConsumer
	MessageHandler *adapters.MessageHandler
}

func NewContainer(config *Config) (*Container, error) {
	storageClient, err := sharedstorage.NewClient(sharedstorage.Config{
		Region:       config.AWSRegion,
		AccessKey:    config.S3AccessKey,
		SecretKey:    config.S3SecretKey,
		Endpoint:     config.S3Endpoint,
		UsePathStyle: config.S3UsePathStyle,
	})
	if err != nil {
		return nil, err
	}

	consumer, err := messaging.NewSQSConsumer(config.AWSRegion, config.SQSQueueURL)
	if err != nil {
		return nil, err
	}

	videoRepo := adapters.NewVideoRepository()
	storageRepo := adapters.NewStorageRepository(storageClient)
	processing := services.NewMP4VideoProcessingService()
	notification := services.NewNotificationService(consumer, config.StateMachineQueue)

	uc := usecases.NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processing,
		notification,
		config.RawBucket,
		config.ProcessedBucket,
		config.MaxSeconds,
	)

	handler := adapters.NewMessageHandler(uc)

	return &Container{Config: config, Consumer: consumer, MessageHandler: handler}, nil
}
