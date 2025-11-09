// internal/infrastructure/container.go
package infrastructure

import (
	"gossipopenclose/internal/adapters"
	"gossipopenclose/internal/application/services"
	"gossipopenclose/internal/application/usecases"
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
		Region:          config.AWSRegion,
		AccessKey:       config.S3AccessKey,
		SecretKey:       config.S3SecretKey,
		SessionToken:    config.S3SessionToken,
		Endpoint:        config.S3Endpoint,
		UsePathStyle:    config.S3UsePathStyle,
		AnonymousAccess: config.S3AnonymousAccess,
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
	processing := services.NewOpenCloseVideoProcessingService()
	notificationService := services.NewNotificationService(consumer, "states_machine_queue")

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

	handler := adapters.NewMessageHandler(uc)
	return &Container{Config: config, Consumer: consumer, MessageHandler: handler}, nil
}
