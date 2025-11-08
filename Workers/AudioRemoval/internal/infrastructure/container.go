package infrastructure

import (
	"audioremoval/internal/adapters"
	"audioremoval/internal/application/services"
	"audioremoval/internal/application/usecases"
	"../../shared/messaging"

	sharedstorage "shared/storage"
)

type Container struct {
	Config         *Config
	Consumer       *messaging.SQSConsumer
	MessageHandler *adapters.MessageHandler
	ProcessVideoUC *usecases.ProcessVideoUseCase
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
	processingService := services.NewMP4VideoProcessingService()
	notificationService := services.NewNotificationService(consumer, config.StateMachineQueue)

	processVideoUC := usecases.NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processingService,
		notificationService,
		config.RawBucket,
		config.ProcessedBucket,
	)

	messageHandler := adapters.NewMessageHandler(processVideoUC)

	return &Container{
		Config:         config,
		Consumer:       consumer,
		MessageHandler: messageHandler,
		ProcessVideoUC: processVideoUC,
	}, nil
}
