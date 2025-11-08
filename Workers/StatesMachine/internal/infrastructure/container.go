package infrastructure

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"shared/messaging"
	"statesmachine/internal/adapters"
	"statesmachine/internal/application/usecases"
)

type Container struct {
	Consumer       *messaging.SQSConsumer
	MessageHandler *adapters.MessageHandler
	DB             *gorm.DB
}

func NewContainer(config *Config) (*Container, error) {
	consumer, err := messaging.NewSQSConsumer(config.AWSRegion, config.SQSQueueURL)
	if err != nil {
		return nil, err
	}

	// Database connection
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	videoRepo := adapters.NewPostgresVideoRepository(db)
	publisher := NewSQSPublisherAdapter(consumer)
	orchestrateUC := usecases.NewOrchestrateVideoUseCase(
		videoRepo,
		publisher,
		config.TrimVideoQueue,
		config.EditVideoQueue,
		config.AudioRemovalQueue,
		config.WatermarkingQueue,
		config.GossipQueue,
		config.MaxRetries,
		config.RetryDelayMinutes,
		config.ProcessedVideoURL,
	)
	messageHandler := adapters.NewMessageHandler(orchestrateUC)

	return &Container{
		Consumer:       consumer,
		MessageHandler: messageHandler,
		DB:             db,
	}, nil
}
