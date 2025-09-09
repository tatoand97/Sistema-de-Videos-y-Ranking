package infrastructure

import (
	"trimvideo/internal/adapters"
	"trimvideo/internal/application/services"
	"trimvideo/internal/application/usecases"
)

type Container struct {
	Config         *Config
	Consumer       *adapters.RabbitMQConsumer
	Publisher      *adapters.RabbitMQPublisher
	MessageHandler *adapters.MessageHandler
}

func NewContainer(config *Config) (*Container, error) {
	storage, err := adapters.NewMinIOStorage(config.MinIOEndpoint, config.MinIOAccessKey, config.MinIOSecretKey)
	if err != nil { return nil, err }

	publisher, err := adapters.NewRabbitMQPublisher(config.RabbitMQURL)
	if err != nil { return nil, err }

	videoRepo := adapters.NewVideoRepository()
	storageRepo := adapters.NewStorageRepository(storage)
	processing := services.NewMP4VideoProcessingService()
	notification := services.NewNotificationService(publisher, config.StateMachineQueue)

	uc := usecases.NewProcessVideoUseCase(
		videoRepo,
		storageRepo,
		processing,
		notification,
		config.RawBucket,
		config.ProcessedBucket,
		config.MaxSeconds,
	)

	consumer, err := adapters.NewRabbitMQConsumer(config.RabbitMQURL, config.MaxRetries, config.QueueMaxLength)
	if err != nil { return nil, err }
	handler := adapters.NewMessageHandler(uc)

	return &Container{Config: config, Consumer: consumer, Publisher: publisher, MessageHandler: handler}, nil
}
