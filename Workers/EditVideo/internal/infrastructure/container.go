package infrastructure

import (
	"editvideo/internal/adapters"
	"editvideo/internal/application/services"
	"editvideo/internal/application/usecases"

	sharedstorage "shared/storage"
)

type Container struct {
	Config         *Config
	Consumer       *adapters.RabbitMQConsumer
	Publisher      *adapters.RabbitMQPublisher
	MessageHandler *adapters.MessageHandler
}

func NewContainer(config *Config) (*Container, error) {
	storageClient, err := sharedstorage.NewClient(sharedstorage.Config{
		Region:       config.S3Region,
		AccessKey:    config.S3AccessKey,
		SecretKey:    config.S3SecretKey,
		Endpoint:     config.S3Endpoint,
		UsePathStyle: config.S3UsePathStyle,
	})
	if err != nil {
		return nil, err
	}

	publisher, err := adapters.NewRabbitMQPublisher(config.RabbitMQURL)
	if err != nil {
		return nil, err
	}

	videoRepo := adapters.NewVideoRepository()
	storageRepo := adapters.NewStorageRepository(storageClient)
	processing := services.NewMP4VideoProcessingService()
	notification := services.NewNotificationService(publisher, config.StateMachineQueue)

	uc := usecases.NewEditVideoUseCase(
		videoRepo,
		storageRepo,
		processing,
		notification,
		config.RawBucket,
		config.ProcessedBucket,
		config.MaxSeconds,
	)

	consumer, err := adapters.NewRabbitMQConsumer(config.RabbitMQURL, config.MaxRetries, config.QueueMaxLength)
	if err != nil {
		return nil, err
	}
	handler := adapters.NewMessageHandler(uc)

	return &Container{Config: config, Consumer: consumer, Publisher: publisher, MessageHandler: handler}, nil
}
