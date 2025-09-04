package infrastructure

import (
	"statesmachine/internal/adapters"
	"statesmachine/internal/application/usecases"
)

type Container struct {
	Consumer       *adapters.RabbitMQConsumer
	MessageHandler *adapters.MessageHandler
}

func NewContainer(config *Config) (*Container, error) {
	consumer, err := adapters.NewRabbitMQConsumer(config.RabbitMQURL)
	if err != nil { return nil, err }

	videoRepo := adapters.NewMockVideoRepository()
	orchestrateUC := usecases.NewOrchestrateVideoUseCase(videoRepo, consumer)
	messageHandler := adapters.NewMessageHandler(orchestrateUC)

	return &Container{
		Consumer:       consumer,
		MessageHandler: messageHandler,
	}, nil
}