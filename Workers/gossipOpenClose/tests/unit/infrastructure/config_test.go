package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		def      int
		envValue string
		expected int
	}{
		{"valid int", "TEST_INT", 10, "25", 25},
		{"invalid int", "TEST_INT", 10, "invalid", 10},
		{"empty value", "TEST_INT", 10, "", 10},
		{"zero value", "TEST_INT", 10, "0", 0},
		{"negative value", "TEST_INT", 10, "-5", -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}
			
			result := getEnvInt(tt.key, tt.def)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvFloat(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		def      float64
		envValue string
		expected float64
	}{
		{"valid float", "TEST_FLOAT", 2.5, "3.7", 3.7},
		{"invalid float", "TEST_FLOAT", 2.5, "invalid", 2.5},
		{"empty value", "TEST_FLOAT", 2.5, "", 2.5},
		{"zero value", "TEST_FLOAT", 2.5, "0.0", 0.0},
		{"negative value", "TEST_FLOAT", 2.5, "-1.5", -1.5},
		{"integer as float", "TEST_FLOAT", 2.5, "5", 5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}
			
			result := getEnvFloat(tt.key, tt.def)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	clearEnvVars()
	
	config := LoadConfig()
	
	assert.NotNil(t, config)
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

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	envVars := map[string]string{
		"RABBITMQ_URL":           "amqp://localhost:5672",
		"MINIO_ENDPOINT":         "localhost:9000",
		"MINIO_ACCESS_KEY":       "minioadmin",
		"MINIO_SECRET_KEY":       "minioadmin",
		"MINIO_BUCKET_RAW":       "raw-videos",
		"MINIO_BUCKET_PROCESSED": "processed-videos",
		"QUEUE_NAME":             "gossip-queue",
		"MAX_RETRIES":            "3",
		"QUEUE_MAX_LENGTH":       "500",
		"MAX_SECONDS":            "60",
		"INTRO_SECONDS":          "3.0",
		"OUTRO_SECONDS":          "2.0",
		"TARGET_WIDTH":           "1920",
		"TARGET_HEIGHT":          "1080",
		"FPS":                    "60",
		"LOGO_PATH":              "/custom/logo.png",
	}
	
	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer clearEnvVars()
	
	config := LoadConfig()
	
	assert.Equal(t, "amqp://localhost:5672", config.RabbitMQURL)
	assert.Equal(t, "localhost:9000", config.MinIOEndpoint)
	assert.Equal(t, "minioadmin", config.MinIOAccessKey)
	assert.Equal(t, "minioadmin", config.MinIOSecretKey)
	assert.Equal(t, "raw-videos", config.RawBucket)
	assert.Equal(t, "processed-videos", config.ProcessedBucket)
	assert.Equal(t, "gossip-queue", config.QueueName)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 500, config.QueueMaxLength)
	assert.Equal(t, 60, config.MaxSeconds)
	assert.Equal(t, 3.0, config.IntroSeconds)
	assert.Equal(t, 2.0, config.OutroSeconds)
	assert.Equal(t, 1920, config.TargetWidth)
	assert.Equal(t, 1080, config.TargetHeight)
	assert.Equal(t, 60, config.FPS)
	assert.Equal(t, "/custom/logo.png", config.LogoPath)
}

func TestLoadConfig_InvalidValues(t *testing.T) {
	clearEnvVars()
	os.Setenv("MAX_RETRIES", "invalid")
	os.Setenv("INTRO_SECONDS", "not-a-float")
	os.Setenv("TARGET_WIDTH", "abc")
	defer clearEnvVars()
	
	config := LoadConfig()
	
	// Should use default values for invalid inputs
	assert.Equal(t, 5, config.MaxRetries)
	assert.Equal(t, 2.5, config.IntroSeconds)
	assert.Equal(t, 1280, config.TargetWidth)
}

func TestLoadConfig_VideoProcessingDefaults(t *testing.T) {
	clearEnvVars()
	
	config := LoadConfig()
	
	// Test video processing specific defaults
	assert.Equal(t, 2.5, config.IntroSeconds)
	assert.Equal(t, 2.5, config.OutroSeconds)
	assert.Equal(t, 1280, config.TargetWidth)
	assert.Equal(t, 720, config.TargetHeight)
	assert.Equal(t, 30, config.FPS)
	assert.Contains(t, config.LogoPath, "nba-logo")
}

func TestConfig_Structure(t *testing.T) {
	config := &Config{
		RabbitMQURL:     "amqp://test:test@localhost:5672/",
		MinIOEndpoint:   "localhost:9000",
		MinIOAccessKey:  "test-key",
		MinIOSecretKey:  "test-secret",
		RawBucket:       "raw",
		ProcessedBucket: "processed",
		QueueName:       "test-queue",
		MaxRetries:      3,
		QueueMaxLength:  100,
		MaxSeconds:      30,
		IntroSeconds:    2.5,
		OutroSeconds:    2.5,
		TargetWidth:     1920,
		TargetHeight:    1080,
		FPS:             30,
		LogoPath:        "/path/to/logo.png",
	}
	
	assert.Equal(t, "amqp://test:test@localhost:5672/", config.RabbitMQURL)
	assert.Equal(t, "localhost:9000", config.MinIOEndpoint)
	assert.Equal(t, "test-key", config.MinIOAccessKey)
	assert.Equal(t, "test-secret", config.MinIOSecretKey)
	assert.Equal(t, "raw", config.RawBucket)
	assert.Equal(t, "processed", config.ProcessedBucket)
	assert.Equal(t, "test-queue", config.QueueName)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 100, config.QueueMaxLength)
	assert.Equal(t, 30, config.MaxSeconds)
	assert.Equal(t, 2.5, config.IntroSeconds)
	assert.Equal(t, 2.5, config.OutroSeconds)
	assert.Equal(t, 1920, config.TargetWidth)
	assert.Equal(t, 1080, config.TargetHeight)
	assert.Equal(t, 30, config.FPS)
	assert.Equal(t, "/path/to/logo.png", config.LogoPath)
}

func clearEnvVars() {
	envVars := []string{
		"RABBITMQ_URL", "MINIO_ENDPOINT", "MINIO_ACCESS_KEY", "MINIO_SECRET_KEY",
		"MINIO_BUCKET_RAW", "MINIO_BUCKET_PROCESSED", "QUEUE_NAME",
		"MAX_RETRIES", "QUEUE_MAX_LENGTH", "MAX_SECONDS",
		"INTRO_SECONDS", "OUTRO_SECONDS", "TARGET_WIDTH", "TARGET_HEIGHT",
		"FPS", "LOGO_PATH",
	}
	
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}