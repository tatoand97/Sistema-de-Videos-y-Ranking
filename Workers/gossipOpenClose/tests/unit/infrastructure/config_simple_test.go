package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear environment variables
	envVars := []string{
		"RABBITMQ_URL", "MINIO_ENDPOINT", "MINIO_ACCESS_KEY", "MINIO_SECRET_KEY",
		"MINIO_BUCKET_RAW", "MINIO_BUCKET_PROCESSED", "QUEUE_NAME",
		"INTRO_SECONDS", "OUTRO_SECONDS", "TARGET_WIDTH", "TARGET_HEIGHT", "FPS",
		"LOGO_PATH", "MAX_RETRIES", "QUEUE_MAX_LENGTH", "MAX_SECONDS",
	}
	
	for _, env := range envVars {
		os.Unsetenv(env)
	}
	
	config := LoadConfig()
	
	// Test default values
	assert.Equal(t, 5, config.MaxRetries)
	assert.Equal(t, 1000, config.QueueMaxLength)
	assert.Equal(t, 30, config.MaxSeconds)
	assert.Equal(t, 2.5, config.IntroSeconds)
	assert.Equal(t, 2.5, config.OutroSeconds)
	assert.Equal(t, 1280, config.TargetWidth)
	assert.Equal(t, 720, config.TargetHeight)
	assert.Equal(t, 30, config.FPS)
	assert.Equal(t, "./assets/nba-logo-removebg-preview.png", config.LogoPath)
}

func TestLoadConfig_EnvironmentOverrides(t *testing.T) {
	// Set environment variables
	testEnvVars := map[string]string{
		"RABBITMQ_URL":           "amqp://test:test@test-rabbit:5672/",
		"MINIO_ENDPOINT":         "test-minio:9001",
		"MINIO_ACCESS_KEY":       "test-access",
		"MINIO_SECRET_KEY":       "test-secret",
		"MINIO_BUCKET_RAW":       "test-raw",
		"MINIO_BUCKET_PROCESSED": "test-processed",
		"QUEUE_NAME":             "test-queue",
		"INTRO_SECONDS":          "3.0",
		"OUTRO_SECONDS":          "1.5",
		"TARGET_WIDTH":           "1920",
		"TARGET_HEIGHT":          "1080",
		"FPS":                    "25",
		"LOGO_PATH":              "/custom/logo.png",
		"MAX_RETRIES":            "10",
		"QUEUE_MAX_LENGTH":       "2000",
		"MAX_SECONDS":            "60",
	}
	
	for key, value := range testEnvVars {
		t.Setenv(key, value)
	}
	
	config := LoadConfig()
	
	// Test overridden values
	assert.Equal(t, "amqp://test:test@test-rabbit:5672/", config.RabbitMQURL)
	assert.Equal(t, "test-minio:9001", config.MinIOEndpoint)
	assert.Equal(t, "test-access", config.MinIOAccessKey)
	assert.Equal(t, "test-secret", config.MinIOSecretKey)
	assert.Equal(t, "test-raw", config.RawBucket)
	assert.Equal(t, "test-processed", config.ProcessedBucket)
	assert.Equal(t, "test-queue", config.QueueName)
	assert.Equal(t, 3.0, config.IntroSeconds)
	assert.Equal(t, 1.5, config.OutroSeconds)
	assert.Equal(t, 1920, config.TargetWidth)
	assert.Equal(t, 1080, config.TargetHeight)
	assert.Equal(t, 25, config.FPS)
	assert.Equal(t, "/custom/logo.png", config.LogoPath)
	assert.Equal(t, 10, config.MaxRetries)
	assert.Equal(t, 2000, config.QueueMaxLength)
	assert.Equal(t, 60, config.MaxSeconds)
}

func TestGetEnvInt_ValidValues(t *testing.T) {
	t.Setenv("TEST_INT", "42")
	
	result := getEnvInt("TEST_INT", 10)
	assert.Equal(t, 42, result)
}

func TestGetEnvInt_InvalidValues(t *testing.T) {
	t.Setenv("INVALID_INT", "not_a_number")
	
	result := getEnvInt("INVALID_INT", 25)
	assert.Equal(t, 25, result) // Should return default
}

func TestGetEnvFloat_ValidValues(t *testing.T) {
	t.Setenv("TEST_FLOAT", "3.14")
	
	result := getEnvFloat("TEST_FLOAT", 1.0)
	assert.Equal(t, 3.14, result)
}

func TestGetEnvFloat_InvalidValues(t *testing.T) {
	t.Setenv("INVALID_FLOAT", "not_a_float")
	
	result := getEnvFloat("INVALID_FLOAT", 2.5)
	assert.Equal(t, 2.5, result) // Should return default
}