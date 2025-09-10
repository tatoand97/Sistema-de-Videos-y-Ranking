package infrastructure_test

import (
	"audioremoval/internal/infrastructure"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_WithDefaults(t *testing.T) {
	// Clear environment variables
	os.Clearenv()
	
	config := infrastructure.LoadConfig()
	
	assert.NotNil(t, config)
	assert.NotEmpty(t, config.RabbitMQURL)
	assert.NotEmpty(t, config.MinioEndpoint)
	assert.NotEmpty(t, config.RawBucket)
	assert.NotEmpty(t, config.ProcessedBucket)
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("RABBITMQ_URL", "amqp://test:5672")
	os.Setenv("MINIO_ENDPOINT", "test-minio:9000")
	os.Setenv("MINIO_ACCESS_KEY", "test-access")
	os.Setenv("MINIO_SECRET_KEY", "test-secret")
	os.Setenv("RAW_BUCKET", "test-raw")
	os.Setenv("PROCESSED_BUCKET", "test-processed")
	os.Setenv("QUEUE_NAME", "test-queue")
	os.Setenv("STATE_QUEUE", "test-state-queue")
	
	defer os.Clearenv()
	
	config := infrastructure.LoadConfig()
	
	assert.Equal(t, "amqp://test:5672", config.RabbitMQURL)
	assert.Equal(t, "test-minio:9000", config.MinioEndpoint)
	assert.Equal(t, "test-access", config.MinioAccessKey)
	assert.Equal(t, "test-secret", config.MinioSecretKey)
	assert.Equal(t, "test-raw", config.RawBucket)
	assert.Equal(t, "test-processed", config.ProcessedBucket)
	assert.Equal(t, "test-queue", config.QueueName)
	assert.Equal(t, "test-state-queue", config.StateQueue)
}

func TestConfig_Validation(t *testing.T) {
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
	
	assert.NotEmpty(t, config.RabbitMQURL)
	assert.NotEmpty(t, config.MinioEndpoint)
	assert.NotEmpty(t, config.MinioAccessKey)
	assert.NotEmpty(t, config.MinioSecretKey)
	assert.NotEmpty(t, config.RawBucket)
	assert.NotEmpty(t, config.ProcessedBucket)
	assert.NotEmpty(t, config.QueueName)
	assert.NotEmpty(t, config.StateQueue)
}

func TestConfig_EmptyValues(t *testing.T) {
	os.Clearenv()
	os.Setenv("RABBITMQ_URL", "")
	os.Setenv("MINIO_ENDPOINT", "")
	
	defer os.Clearenv()
	
	config := infrastructure.LoadConfig()
	
	// Should use defaults when empty
	assert.NotEmpty(t, config.RabbitMQURL)
	assert.NotEmpty(t, config.MinioEndpoint)
}