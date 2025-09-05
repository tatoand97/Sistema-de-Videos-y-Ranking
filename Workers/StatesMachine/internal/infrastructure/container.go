package infrastructure

import (
	"statesmachine/internal/adapters"
	"statesmachine/internal/application/usecases"
)

type Container struct {
	Consumer       *adapters.RabbitMQConsumer
	Publisher      *adapters.RabbitMQPublisher
	MessageHandler *adapters.MessageHandler
}

func NewContainer(config *Config) (*Container, error) {
	consumer, err := adapters.NewRabbitMQConsumer(config.RabbitMQURL)
	if err != nil { return nil, err }

	publisher, err := adapters.NewRabbitMQPublisher(config.RabbitMQURL)
	if err != nil { return nil, err }

	videoRepo := adapters.NewMockVideoRepository()
	orchestrateUC := usecases.NewOrchestrateVideoUseCase(videoRepo, publisher, config.EditVideoQueue, config.AudioRemovalQueue)
	messageHandler := adapters.NewMessageHandler(orchestrateUC)

	return &Container{
		Consumer:       consumer,
		Publisher:      publisher,
		MessageHandler: messageHandler,
	}, nil
}