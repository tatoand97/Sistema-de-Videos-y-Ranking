package infrastructure

import (
	"watermarking/internal/adapters"
	"watermarking/internal/application/services"
	"watermarking/internal/application/usecases"
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
	processing := services.NewMP4VideoProcessingService()

	uc := usecases.NewWatermarkingUseCase(
		videoRepo,
		storageRepo,
		processing,
		config.RawBucket,
		config.ProcessedBucket,
		config.MaxSeconds,
	)

	consumer, err := adapters.NewRabbitMQConsumer(config.RabbitMQURL, config.MaxRetries, config.QueueMaxLength)
	if err != nil { return nil, err }
	handler := adapters.NewMessageHandler(uc)

	return &Container{Config: config, Consumer: consumer, MessageHandler: handler}, nil
}
