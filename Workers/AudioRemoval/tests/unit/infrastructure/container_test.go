package infrastructure_test

import (
	"audioremoval/internal/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer(t *testing.T) {
	config := &infrastructure.Config{
		RabbitMQURL:     "amqp://localhost:5672",
		MinioEndpoint:   "localhost:9000",
		MinioAccessKey:  "access",
		MinioSecretKey:  "secret",
		RawBucket:       "raw",
		ProcessedBucket: "processed",
		QueueName:       "queue",
		StateQueue:      "state",
	}
	
	container := infrastructure.NewContainer(config)
	
	assert.NotNil(t, container)
}

func TestContainer_GetProcessVideoUseCase(t *testing.T) {
	config := &infrastructure.Config{
		RabbitMQURL:     "amqp://localhost:5672",
		MinioEndpoint:   "localhost:9000",
		MinioAccessKey:  "access",
		MinioSecretKey:  "secret",
		RawBucket:       "raw",
		ProcessedBucket: "processed",
		QueueName:       "queue",
		StateQueue:      "state",
	}
	
	container := infrastructure.NewContainer(config)
	useCase := container.GetProcessVideoUseCase()
	
	assert.NotNil(t, useCase)
}

func TestContainer_GetMessageHandler(t *testing.T) {
	config := &infrastructure.Config{
		RabbitMQURL:     "amqp://localhost:5672",
		MinioEndpoint:   "localhost:9000",
		MinioAccessKey:  "access",
		MinioSecretKey:  "secret",
		RawBucket:       "raw",
		ProcessedBucket: "processed",
		QueueName:       "queue",
		StateQueue:      "state",
	}
	
	container := infrastructure.NewContainer(config)
	handler := container.GetMessageHandler()
	
	assert.NotNil(t, handler)
}

func TestContainer_GetRabbitMQConsumer(t *testing.T) {
	config := &infrastructure.Config{
		RabbitMQURL:     "amqp://localhost:5672",
		MinioEndpoint:   "localhost:9000",
		MinioAccessKey:  "access",
		MinioSecretKey:  "secret",
		RawBucket:       "raw",
		ProcessedBucket: "processed",
		QueueName:       "queue",
		StateQueue:      "state",
	}
	
	container := infrastructure.NewContainer(config)
	consumer := container.GetRabbitMQConsumer()
	
	assert.NotNil(t, consumer)
}

func TestContainer_WithNilConfig(t *testing.T) {
	container := infrastructure.NewContainer(nil)
	
	assert.NotNil(t, container)
}

func TestContainer_MultipleCallsSameInstance(t *testing.T) {
	config := &infrastructure.Config{
		RabbitMQURL:     "amqp://localhost:5672",
		MinioEndpoint:   "localhost:9000",
		MinioAccessKey:  "access",
		MinioSecretKey:  "secret",
		RawBucket:       "raw",
		ProcessedBucket: "processed",
		QueueName:       "queue",
		StateQueue:      "state",
	}
	
	container := infrastructure.NewContainer(config)
	
	useCase1 := container.GetProcessVideoUseCase()
	useCase2 := container.GetProcessVideoUseCase()
	
	// Should return same instance (singleton pattern)
	assert.Equal(t, useCase1, useCase2)
}