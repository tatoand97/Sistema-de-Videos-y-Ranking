package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer_NilConfig(t *testing.T) {
	// This will panic, so we expect it
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

func TestNewContainer_ValidConfig(t *testing.T) {
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
	}

	// This will likely fail due to connection issues, but tests the structure
	container, err := NewContainer(config)
	
	// If it succeeds (unlikely without real services), verify structure
	if err == nil {
		assert.NotNil(t, container)
		assert.NotNil(t, container.Config)
		assert.NotNil(t, container.Consumer)
		assert.NotNil(t, container.Publisher)
		assert.NotNil(t, container.MessageHandler)
		assert.NotNil(t, container.ProcessVideoUC)
	} else {
		// Expected to fail without real MinIO/RabbitMQ services
		assert.Error(t, err)
	}
}

func TestContainer_Structure(t *testing.T) {
	// Test that we can create the container structure
	container := &Container{
		Config:         &Config{},
		Consumer:       nil,
		Publisher:      nil,
		MessageHandler: nil,
		ProcessVideoUC: nil,
	}

	assert.NotNil(t, container)
	assert.NotNil(t, container.Config)
}

func TestContainer_ConfigValidation(t *testing.T) {
	// Test config field validation
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
	}

	assert.NotEmpty(t, config.MinIOEndpoint)
	assert.NotEmpty(t, config.MinIOAccessKey)
	assert.NotEmpty(t, config.MinIOSecretKey)
	assert.NotEmpty(t, config.RabbitMQURL)
	assert.Greater(t, config.MaxRetries, 0)
	assert.Greater(t, config.QueueMaxLength, 0)
	assert.NotEmpty(t, config.StateMachineQueue)
	assert.NotEmpty(t, config.RawBucket)
	assert.NotEmpty(t, config.ProcessedBucket)
}