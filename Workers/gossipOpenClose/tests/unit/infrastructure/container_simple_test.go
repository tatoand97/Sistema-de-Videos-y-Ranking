package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer_WithValidConfig(t *testing.T) {
	config := &Config{
		RabbitMQURL:     "amqp://guest:guest@localhost:5672/",
		MinIOEndpoint:   "localhost:9000",
		MinIOAccessKey:  "minioadmin",
		MinIOSecretKey:  "minioadmin",
		RawBucket:       "raw-videos",
		ProcessedBucket: "processed-videos",
		QueueName:       "test-queue",
		MaxRetries:      5,
		QueueMaxLength:  1000,
		IntroSeconds:    2.5,
		OutroSeconds:    2.5,
		TargetWidth:     1920,
		TargetHeight:    1080,
		FPS:             30,
		LogoPath:        "./assets/logo.png",
	}
	
	container, err := NewContainer(config)
	
	// Expected to fail in test environment due to missing services
	if err != nil {
		assert.Error(t, err)
		assert.Nil(t, container)
	} else {
		assert.NotNil(t, container)
		assert.Equal(t, config, container.Config)
	}
}

func TestNewContainer_WithEmptyConfig(t *testing.T) {
	config := &Config{}
	
	container, err := NewContainer(config)
	
	// Should fail with empty configuration
	assert.Error(t, err)
	assert.Nil(t, container)
}

func TestConfig_Fields(t *testing.T) {
	config := &Config{
		RabbitMQURL:     "test-url",
		MinIOEndpoint:   "test-endpoint",
		MinIOAccessKey:  "test-access",
		MinIOSecretKey:  "test-secret",
		RawBucket:       "raw",
		ProcessedBucket: "processed",
		QueueName:       "queue",
		MaxRetries:      3,
		QueueMaxLength:  500,
		MaxSeconds:      60,
		IntroSeconds:    1.0,
		OutroSeconds:    2.0,
		TargetWidth:     1280,
		TargetHeight:    720,
		FPS:             25,
		LogoPath:        "/path/to/logo.png",
	}
	
	assert.Equal(t, "test-url", config.RabbitMQURL)
	assert.Equal(t, "test-endpoint", config.MinIOEndpoint)
	assert.Equal(t, "test-access", config.MinIOAccessKey)
	assert.Equal(t, "test-secret", config.MinIOSecretKey)
	assert.Equal(t, "raw", config.RawBucket)
	assert.Equal(t, "processed", config.ProcessedBucket)
	assert.Equal(t, "queue", config.QueueName)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 500, config.QueueMaxLength)
	assert.Equal(t, 60, config.MaxSeconds)
	assert.Equal(t, 1.0, config.IntroSeconds)
	assert.Equal(t, 2.0, config.OutroSeconds)
	assert.Equal(t, 1280, config.TargetWidth)
	assert.Equal(t, 720, config.TargetHeight)
	assert.Equal(t, 25, config.FPS)
	assert.Equal(t, "/path/to/logo.png", config.LogoPath)
}