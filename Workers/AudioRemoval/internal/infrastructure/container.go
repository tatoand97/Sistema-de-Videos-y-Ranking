package infrastructure

import (
	"audioremoval/internal/adapters"
	"audioremoval/internal/application/services"
	"audioremoval/internal/application/usecases"
)

type Container struct {
	Config         *Config
	Consumer       *adapters.RabbitMQConsumer
	Publisher      *adapters.RabbitMQPublisher
	MessageHandler *adapters.MessageHandler
	ProcessVideoUC *usecases.ProcessVideoUseCase
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

	publisher, err := adapters.NewRabbitMQPublisher(config.RabbitMQURL)
	if err != nil {
		return nil, err
	}

	consumer, err := adapters.NewRabbitMQConsumer(config.RabbitMQURL, config.MaxRetries, config.QueueMaxLength)
	if err != nil {
		return nil, err
	}

	videoRepo := adapters.NewVideoRepository()
	storageRepo := adapters.NewStorageRepository(storage)
	processingService := services.NewMP4VideoProcessingService()
	notificationService := services.NewNotificationService(publisher, config.StateMachineQueue)

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
		Publisher:      publisher,
		MessageHandler: messageHandler,
		ProcessVideoUC: processVideoUC,
	}, nil
}