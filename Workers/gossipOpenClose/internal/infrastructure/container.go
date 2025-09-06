// internal/infrastructure/container.go
package infrastructure

import (
	"gossipopenclose/internal/adapters"
	"gossipopenclose/internal/application/services"
	"gossipopenclose/internal/application/usecases"
)

type Container struct {
	Config         *Config
	Consumer       *adapters.RabbitMQConsumer
	MessageHandler *adapters.MessageHandler
}

func NewContainer(config *Config) (*Container, error) {
	storage, err := adapters.NewMinIOStorage(config.MinIOEndpoint, config.MinIOAccessKey, config.MinIOSecretKey)
	if err != nil { return nil, err }

	videoRepo := adapters.NewVideoRepository()
	storageRepo := adapters.NewStorageRepository(storage)

	processing := services.NewOpenCloseVideoProcessingService()

	publisher, err := adapters.NewRabbitMQPublisher(config.RabbitMQURL)
	if err != nil { return nil, err }

	notificationService := services.NewNotificationService(publisher, "states_machine_queue")

	uc := usecases.NewOpenCloseUseCase(
		videoRepo,
		storageRepo,
		processing,
		notificationService,
		config.RawBucket,
		config.ProcessedBucket,
		config.LogoPath,
		config.IntroSeconds,
		config.OutroSeconds,
		config.TargetWidth,
		config.TargetHeight,
		config.FPS,
	)

	consumer, err := adapters.NewRabbitMQConsumer(config.RabbitMQURL, config.MaxRetries, config.QueueMaxLength)
	if err != nil { return nil, err }

	handler := adapters.NewMessageHandler(uc)
	return &Container{Config: config, Consumer: consumer, MessageHandler: handler}, nil
}