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

func TestNewContainer_InvalidS3Config(t *testing.T) {
	config := &Config{
		S3Region:    "",
		RabbitMQURL: "amqp://localhost:5672",
	}

	container, err := NewContainer(config)
	assert.Error(t, err)
	assert.Nil(t, container)
}

func TestNewContainer_InvalidRabbitMQConfig(t *testing.T) {
	config := &Config{
		S3Region:    "us-east-1",
		RabbitMQURL: "invalid-url",
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
		S3Region:          "us-east-1",
		S3Endpoint:        "https://s3.us-east-1.amazonaws.com",
		RabbitMQURL:       "amqp://guest:guest@localhost:5672/",
		MaxRetries:        3,
		QueueMaxLength:    1000,
		StateMachineQueue: "state_machine_queue",
		RawBucket:         "raw-videos",
		ProcessedBucket:   "processed-videos",
		MaxSeconds:        30,
	}

	assert.NotEmpty(t, config.S3Region)
	assert.NotEmpty(t, config.RabbitMQURL)
	assert.Greater(t, config.MaxRetries, 0)
	assert.Greater(t, config.MaxSeconds, 0)
}
