package infrastructure_test

import (
	"audioremoval/internal/infrastructure"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_WithDefaults(t *testing.T) {
	os.Clearenv()

	config := infrastructure.LoadConfig()

	assert.NotNil(t, config)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 1000, config.QueueMaxLength)
	assert.False(t, config.S3UsePathStyle)
	assert.Empty(t, config.S3Region)
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	os.Clearenv()
	t.Cleanup(func() { os.Clearenv() })

	os.Setenv("RABBITMQ_URL", "amqp://test:5672")
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("S3_ENDPOINT", "https://s3.us-east-1.amazonaws.com")
	os.Setenv("AWS_ACCESS_KEY_ID", "test-access")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test-secret")
	os.Setenv("S3_USE_PATH_STYLE", "true")
	os.Setenv("RAW_BUCKET", "test-raw")
	os.Setenv("PROCESSED_BUCKET", "test-processed")
	os.Setenv("QUEUE_NAME", "test-queue")
	os.Setenv("STATE_MACHINE_QUEUE", "test-state-queue")
	os.Setenv("MAX_RETRIES", "7")
	os.Setenv("QUEUE_MAX_LENGTH", "1500")

	config := infrastructure.LoadConfig()

	assert.Equal(t, "amqp://test:5672", config.RabbitMQURL)
	assert.Equal(t, "us-east-1", config.S3Region)
	assert.Equal(t, "https://s3.us-east-1.amazonaws.com", config.S3Endpoint)
	assert.Equal(t, "test-access", config.S3AccessKey)
	assert.Equal(t, "test-secret", config.S3SecretKey)
	assert.True(t, config.S3UsePathStyle)
	assert.Equal(t, "test-raw", config.RawBucket)
	assert.Equal(t, "test-processed", config.ProcessedBucket)
	assert.Equal(t, "test-queue", config.QueueName)
	assert.Equal(t, "test-state-queue", config.StateMachineQueue)
	assert.Equal(t, 7, config.MaxRetries)
	assert.Equal(t, 1500, config.QueueMaxLength)
}

func TestConfig_Validation(t *testing.T) {
	config := &infrastructure.Config{
		RabbitMQURL:       "amqp://localhost:5672",
		S3Region:          "us-east-1",
		S3Endpoint:        "https://s3.us-east-1.amazonaws.com",
		S3AccessKey:       "access",
		S3SecretKey:       "secret",
		RawBucket:         "raw",
		ProcessedBucket:   "processed",
		QueueName:         "queue",
		StateMachineQueue: "state",
	}

	assert.NotEmpty(t, config.RabbitMQURL)
	assert.NotEmpty(t, config.S3Region)
	assert.NotEmpty(t, config.S3Endpoint)
	assert.NotEmpty(t, config.S3AccessKey)
	assert.NotEmpty(t, config.S3SecretKey)
	assert.NotEmpty(t, config.RawBucket)
	assert.NotEmpty(t, config.ProcessedBucket)
	assert.NotEmpty(t, config.QueueName)
	assert.NotEmpty(t, config.StateMachineQueue)
}
