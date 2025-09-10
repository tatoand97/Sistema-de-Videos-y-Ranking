package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer_NilConfig(t *testing.T) {
	assert.Panics(t, func() {
		NewContainer(nil)
	})
}

func TestNewContainer_InvalidMinIOConfig(t *testing.T) {
	config := &Config{
		MinIOEndpoint:  "",
		MinIOAccessKey: "",
		MinIOSecretKey: "",
		RabbitMQURL:    "amqp://localhost:5672",
	}

	container, err := NewContainer(config)
	assert.Error(t, err)
	assert.Nil(t, container)
}

func TestNewContainer_InvalidRabbitMQConfig(t *testing.T) {
	config := &Config{
		MinIOEndpoint:  "localhost:9000",
		MinIOAccessKey: "minioadmin",
		MinIOSecretKey: "minioadmin",
		RabbitMQURL:    "invalid-url",
	}

	container, err := NewContainer(config)
	assert.Error(t, err)
	assert.Nil(t, container)
}

func TestContainer_Structure(t *testing.T) {
	container := &Container{
		Config:         &Config{},
		Consumer:       nil,
		Publisher:      nil,
		MessageHandler: nil,
	}

	assert.NotNil(t, container)
	assert.NotNil(t, container.Config)
}

func TestContainer_ConfigValidation(t *testing.T) {
	config := &Config{
		MinIOEndpoint:     "localhost:9000",
		MinIOAccessKey:    "minioadmin",
		MinIOSecretKey:    "minioadmin",
		RabbitMQURL:       "amqp://guest:guest@localhost:5672/",
		MaxRetries:        3,
		QueueMaxLength:    1000,
		StateMachineQueue: "state_machine_queue",
		RawBucket:         "raw-videos",
		ProcessedBucket:   "processed-videos",
		MaxSeconds:        30,
	}

	assert.NotEmpty(t, config.MinIOEndpoint)
	assert.NotEmpty(t, config.RabbitMQURL)
	assert.Greater(t, config.MaxRetries, 0)
	assert.Greater(t, config.MaxSeconds, 0)
}