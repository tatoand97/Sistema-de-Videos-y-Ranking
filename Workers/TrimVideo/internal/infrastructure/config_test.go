package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	clearEnvVars()

	config := LoadConfig()

	assert.NotNil(t, config)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 1000, config.QueueMaxLength)
	assert.Equal(t, 30, config.MaxSeconds)
	assert.Empty(t, config.RabbitMQURL)
	assert.Empty(t, config.S3Region)
	assert.False(t, config.S3UsePathStyle)
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	envVars := map[string]string{
		"RABBITMQ_URL":          "amqp://localhost:5672",
		"AWS_REGION":            "us-east-1",
		"S3_ENDPOINT":           "https://s3.us-east-1.amazonaws.com",
		"AWS_ACCESS_KEY_ID":     "access123",
		"AWS_SECRET_ACCESS_KEY": "secret123",
		"S3_USE_PATH_STYLE":     "true",
		"RAW_BUCKET":            "raw-videos",
		"PROCESSED_BUCKET":      "processed-videos",
		"QUEUE_NAME":            "trim-video",
		"STATE_MACHINE_QUEUE":   "state-machine",
		"MAX_RETRIES":           "5",
		"QUEUE_MAX_LENGTH":      "2000",
		"MAX_SECONDS":           "60",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer clearEnvVars()

	config := LoadConfig()

	assert.Equal(t, "amqp://localhost:5672", config.RabbitMQURL)
	assert.Equal(t, "us-east-1", config.S3Region)
	assert.Equal(t, "https://s3.us-east-1.amazonaws.com", config.S3Endpoint)
	assert.Equal(t, "access123", config.S3AccessKey)
	assert.Equal(t, "secret123", config.S3SecretKey)
	assert.True(t, config.S3UsePathStyle)
	assert.Equal(t, "raw-videos", config.RawBucket)
	assert.Equal(t, "processed-videos", config.ProcessedBucket)
	assert.Equal(t, "trim-video", config.QueueName)
	assert.Equal(t, "state-machine", config.StateMachineQueue)
	assert.Equal(t, 5, config.MaxRetries)
	assert.Equal(t, 2000, config.QueueMaxLength)
	assert.Equal(t, 60, config.MaxSeconds)
}

func TestLoadConfig_InvalidValues(t *testing.T) {
	clearEnvVars()
	os.Setenv("MAX_RETRIES", "invalid")
	os.Setenv("QUEUE_MAX_LENGTH", "not-a-number")
	os.Setenv("MAX_SECONDS", "invalid")
	defer clearEnvVars()

	config := LoadConfig()

	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 1000, config.QueueMaxLength)
	assert.Equal(t, 30, config.MaxSeconds)
}

func TestLoadConfig_ZeroValues(t *testing.T) {
	clearEnvVars()
	os.Setenv("MAX_RETRIES", "0")
	os.Setenv("QUEUE_MAX_LENGTH", "0")
	os.Setenv("MAX_SECONDS", "0")
	defer clearEnvVars()

	config := LoadConfig()

	assert.Equal(t, 0, config.MaxRetries)
	assert.Equal(t, 0, config.QueueMaxLength)
	assert.Equal(t, 0, config.MaxSeconds)
}

func clearEnvVars() {
	envVars := []string{
		"RABBITMQ_URL", "AWS_REGION", "S3_ENDPOINT", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY",
		"S3_USE_PATH_STYLE", "RAW_BUCKET", "PROCESSED_BUCKET", "QUEUE_NAME", "STATE_MACHINE_QUEUE",
		"MAX_RETRIES", "QUEUE_MAX_LENGTH", "MAX_SECONDS",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
