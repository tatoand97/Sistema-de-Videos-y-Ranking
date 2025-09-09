package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Clear environment variables
	clearEnvVars()
	
	config := LoadConfig()
	
	assert.NotNil(t, config)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 1000, config.QueueMaxLength)
	assert.Empty(t, config.RabbitMQURL)
	assert.Empty(t, config.MinIOEndpoint)
	assert.Empty(t, config.MinIOAccessKey)
	assert.Empty(t, config.MinIOSecretKey)
	assert.Empty(t, config.RawBucket)
	assert.Empty(t, config.ProcessedBucket)
	assert.Empty(t, config.QueueName)
	assert.Empty(t, config.StateMachineQueue)
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	envVars := map[string]string{
		"RABBITMQ_URL":         "amqp://localhost:5672",
		"MINIO_ENDPOINT":       "localhost:9000",
		"MINIO_ACCESS_KEY":     "minioadmin",
		"MINIO_SECRET_KEY":     "minioadmin",
		"RAW_BUCKET":           "raw-videos",
		"PROCESSED_BUCKET":     "processed-videos",
		"QUEUE_NAME":           "audio-removal",
		"STATE_MACHINE_QUEUE":  "state-machine",
		"MAX_RETRIES":          "5",
		"QUEUE_MAX_LENGTH":     "2000",
	}
	
	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer clearEnvVars()
	
	config := LoadConfig()
	
	assert.NotNil(t, config)
	assert.Equal(t, "amqp://localhost:5672", config.RabbitMQURL)
	assert.Equal(t, "localhost:9000", config.MinIOEndpoint)
	assert.Equal(t, "minioadmin", config.MinIOAccessKey)
	assert.Equal(t, "minioadmin", config.MinIOSecretKey)
	assert.Equal(t, "raw-videos", config.RawBucket)
	assert.Equal(t, "processed-videos", config.ProcessedBucket)
	assert.Equal(t, "audio-removal", config.QueueName)
	assert.Equal(t, "state-machine", config.StateMachineQueue)
	assert.Equal(t, 5, config.MaxRetries)
	assert.Equal(t, 2000, config.QueueMaxLength)
}

func TestLoadConfig_InvalidMaxRetries(t *testing.T) {
	clearEnvVars()
	os.Setenv("MAX_RETRIES", "invalid")
	defer clearEnvVars()
	
	config := LoadConfig()
	
	// Should use default value when parsing fails
	assert.Equal(t, 3, config.MaxRetries)
}

func TestLoadConfig_InvalidQueueMaxLength(t *testing.T) {
	clearEnvVars()
	os.Setenv("QUEUE_MAX_LENGTH", "not-a-number")
	defer clearEnvVars()
	
	config := LoadConfig()
	
	// Should use default value when parsing fails
	assert.Equal(t, 1000, config.QueueMaxLength)
}

func TestLoadConfig_ZeroValues(t *testing.T) {
	clearEnvVars()
	os.Setenv("MAX_RETRIES", "0")
	os.Setenv("QUEUE_MAX_LENGTH", "0")
	defer clearEnvVars()
	
	config := LoadConfig()
	
	assert.Equal(t, 0, config.MaxRetries)
	assert.Equal(t, 0, config.QueueMaxLength)
}

func TestLoadConfig_NegativeValues(t *testing.T) {
	clearEnvVars()
	os.Setenv("MAX_RETRIES", "-1")
	os.Setenv("QUEUE_MAX_LENGTH", "-100")
	defer clearEnvVars()
	
	config := LoadConfig()
	
	assert.Equal(t, -1, config.MaxRetries)
	assert.Equal(t, -100, config.QueueMaxLength)
}

func TestLoadConfig_LargeValues(t *testing.T) {
	clearEnvVars()
	os.Setenv("MAX_RETRIES", "999999")
	os.Setenv("QUEUE_MAX_LENGTH", "1000000")
	defer clearEnvVars()
	
	config := LoadConfig()
	
	assert.Equal(t, 999999, config.MaxRetries)
	assert.Equal(t, 1000000, config.QueueMaxLength)
}

func TestLoadConfig_EmptyStringValues(t *testing.T) {
	clearEnvVars()
	os.Setenv("MAX_RETRIES", "")
	os.Setenv("QUEUE_MAX_LENGTH", "")
	defer clearEnvVars()
	
	config := LoadConfig()
	
	// Should use default values for empty strings
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 1000, config.QueueMaxLength)
}

func TestLoadConfig_PartialConfiguration(t *testing.T) {
	clearEnvVars()
	os.Setenv("RABBITMQ_URL", "amqp://localhost:5672")
	os.Setenv("MAX_RETRIES", "10")
	// Leave other variables unset
	defer clearEnvVars()
	
	config := LoadConfig()
	
	assert.Equal(t, "amqp://localhost:5672", config.RabbitMQURL)
	assert.Equal(t, 10, config.MaxRetries)
	assert.Equal(t, 1000, config.QueueMaxLength) // default
	assert.Empty(t, config.MinIOEndpoint)
	assert.Empty(t, config.QueueName)
}

func TestConfig_Structure(t *testing.T) {
	config := &Config{
		RabbitMQURL:       "test-url",
		MinIOEndpoint:     "test-endpoint",
		MinIOAccessKey:    "test-key",
		MinIOSecretKey:    "test-secret",
		RawBucket:         "raw",
		ProcessedBucket:   "processed",
		QueueName:         "queue",
		StateMachineQueue: "state",
		MaxRetries:        5,
		QueueMaxLength:    100,
	}
	
	assert.Equal(t, "test-url", config.RabbitMQURL)
	assert.Equal(t, "test-endpoint", config.MinIOEndpoint)
	assert.Equal(t, "test-key", config.MinIOAccessKey)
	assert.Equal(t, "test-secret", config.MinIOSecretKey)
	assert.Equal(t, "raw", config.RawBucket)
	assert.Equal(t, "processed", config.ProcessedBucket)
	assert.Equal(t, "queue", config.QueueName)
	assert.Equal(t, "state", config.StateMachineQueue)
	assert.Equal(t, 5, config.MaxRetries)
	assert.Equal(t, 100, config.QueueMaxLength)
}

// Helper function to clear environment variables
func clearEnvVars() {
	envVars := []string{
		"RABBITMQ_URL",
		"MINIO_ENDPOINT",
		"MINIO_ACCESS_KEY",
		"MINIO_SECRET_KEY",
		"RAW_BUCKET",
		"PROCESSED_BUCKET",
		"QUEUE_NAME",
		"STATE_MACHINE_QUEUE",
		"MAX_RETRIES",
		"QUEUE_MAX_LENGTH",
	}
	
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}