package infrastructure

import (
	"statesmachine/internal/adapters"
	"statesmachine/internal/application/usecases"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Container struct {
	Consumer       *adapters.RabbitMQConsumer
	Publisher      *adapters.RabbitMQPublisher
	MessageHandler *adapters.MessageHandler
	DB             *gorm.DB
}

func NewContainer(config *Config) (*Container, error) {
	consumer, err := adapters.NewRabbitMQConsumer(config.RabbitMQURL)
	if err != nil { return nil, err }

	publisher, err := adapters.NewRabbitMQPublisher(config.RabbitMQURL)
	if err != nil { return nil, err }

	// Database connection
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil { return nil, err }

	videoRepo := adapters.NewPostgresVideoRepository(db)
	orchestrateUC := usecases.NewOrchestrateVideoUseCase(videoRepo, publisher, config.EditVideoQueue, config.AudioRemovalQueue, config.WatermarkingQueue)
	messageHandler := adapters.NewMessageHandler(orchestrateUC)

	return &Container{
		Consumer:       consumer,
		Publisher:      publisher,
		MessageHandler: messageHandler,
		DB:             db,
	}, nil
}