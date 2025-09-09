package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{"with env value", "TEST_KEY", "default", "env_value", "env_value"},
		{"without env value", "MISSING_KEY", "default", "", "default"},
		{"empty env value", "EMPTY_KEY", "default", "", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}
			
			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		expected     int
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
			
			result := getEnvInt(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateRabbitMQURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		expectErr bool
		errMsg    string
	}{
		{"valid amqp url", "amqp://admin:secret@localhost:5672/", false, ""},
		{"valid amqps url", "amqps://admin:secret@localhost:5672/", false, ""},
		{"empty url", "", true, "RABBITMQ_URL is required"},
		{"invalid scheme", "http://localhost:5672/", true, "must use amqp or amqps scheme"},
		{"placeholder credentials", "amqp://user:pass@localhost:5672/", true, "placeholder credentials"},
		{"malformed url", "not-a-url", true, "must use amqp or amqps scheme"},
		{"no credentials", "amqp://localhost:5672/", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRabbitMQURL(tt.url)
			
			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	clearEnvVars()
	
	// Set valid RabbitMQ URL to avoid fatal error
	os.Setenv("RABBITMQ_URL", "amqp://admin:secret@localhost:5672/")
	defer os.Unsetenv("RABBITMQ_URL")
	
	config := LoadConfig()
	
	assert.NotNil(t, config)
	assert.Equal(t, "amqp://admin:secret@localhost:5672/", config.RabbitMQURL)
	assert.Equal(t, "orders", config.QueueName)
	assert.Equal(t, "edit_video_queue", config.EditVideoQueue)
	assert.Equal(t, "audio_removal_queue", config.AudioRemovalQueue)
	assert.Equal(t, "watermarking_queue", config.WatermarkingQueue)
	assert.Contains(t, config.DatabaseURL, "postgres://")
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 5, config.RetryDelayMinutes)
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	envVars := map[string]string{
		"RABBITMQ_URL":         "amqp://test:test@rabbitmq:5672/",
		"QUEUE_NAME":           "test-orders",
		"EDIT_VIDEO_QUEUE":     "test-edit",
		"AUDIO_REMOVAL_QUEUE":  "test-audio",
		"WATERMARKING_QUEUE":   "test-watermark",
		"DATABASE_URL":         "postgres://test:test@db:5432/test",
		"MAX_RETRIES":          "5",
		"RETRY_DELAY_MINUTES":  "10",
	}
	
	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer clearEnvVars()
	
	config := LoadConfig()
	
	assert.Equal(t, "amqp://test:test@rabbitmq:5672/", config.RabbitMQURL)
	assert.Equal(t, "test-orders", config.QueueName)
	assert.Equal(t, "test-edit", config.EditVideoQueue)
	assert.Equal(t, "test-audio", config.AudioRemovalQueue)
	assert.Equal(t, "test-watermark", config.WatermarkingQueue)
	assert.Equal(t, "postgres://test:test@db:5432/test", config.DatabaseURL)
	assert.Equal(t, 5, config.MaxRetries)
	assert.Equal(t, 10, config.RetryDelayMinutes)
}

func TestLoadConfig_InvalidIntValues(t *testing.T) {
	clearEnvVars()
	os.Setenv("RABBITMQ_URL", "amqp://admin:secret@localhost:5672/")
	os.Setenv("MAX_RETRIES", "invalid")
	os.Setenv("RETRY_DELAY_MINUTES", "not-a-number")
	defer clearEnvVars()
	
	config := LoadConfig()
	
	// Should use default values for invalid integers
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 5, config.RetryDelayMinutes)
}

func TestConfig_Structure(t *testing.T) {
	config := &Config{
		RabbitMQURL:         "amqp://test:test@localhost:5672/",
		QueueName:           "test-queue",
		EditVideoQueue:      "edit-queue",
		AudioRemovalQueue:   "audio-queue",
		WatermarkingQueue:   "watermark-queue",
		DatabaseURL:         "postgres://test:test@db:5432/test",
		MaxRetries:          3,
		RetryDelayMinutes:   5,
	}
	
	assert.Equal(t, "amqp://test:test@localhost:5672/", config.RabbitMQURL)
	assert.Equal(t, "test-queue", config.QueueName)
	assert.Equal(t, "edit-queue", config.EditVideoQueue)
	assert.Equal(t, "audio-queue", config.AudioRemovalQueue)
	assert.Equal(t, "watermark-queue", config.WatermarkingQueue)
	assert.Equal(t, "postgres://test:test@db:5432/test", config.DatabaseURL)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 5, config.RetryDelayMinutes)
}

func clearEnvVars() {
	envVars := []string{
		"RABBITMQ_URL", "QUEUE_NAME", "EDIT_VIDEO_QUEUE", "AUDIO_REMOVAL_QUEUE",
		"WATERMARKING_QUEUE", "DATABASE_URL", "MAX_RETRIES", "RETRY_DELAY_MINUTES",
	}
	
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}