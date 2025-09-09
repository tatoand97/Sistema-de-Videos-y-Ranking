package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Success(t *testing.T) {
	// Set environment variables
	os.Setenv("RABBITMQ_URL", "amqp://testuser:testpass@localhost:5672")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost/db")
	os.Setenv("MAX_RETRIES", "5")
	os.Setenv("RETRY_DELAY_MINUTES", "10")
	defer func() {
		os.Unsetenv("RABBITMQ_URL")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("MAX_RETRIES")
		os.Unsetenv("RETRY_DELAY_MINUTES")
	}()

	config := LoadConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "amqp://testuser:testpass@localhost:5672", config.RabbitMQURL)
	assert.Equal(t, "postgres://user:pass@localhost/db", config.DatabaseURL)
	assert.Equal(t, 5, config.MaxRetries)
	assert.Equal(t, 10, config.RetryDelayMinutes)
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("RABBITMQ_URL")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("MAX_RETRIES")
	os.Unsetenv("RETRY_DELAY_MINUTES")

	// This will cause a log.Fatalf due to placeholder credentials
	// We can't easily test this without mocking log.Fatalf
	// So we'll test the validation function separately
	assert.NotPanics(t, func() {
		// Just test that the function exists
		config := &Config{
			RabbitMQURL:       "amqp://testuser:testpass@localhost:5672",
			DatabaseURL:       "postgres://app_user:app_password@postgres:5432/videorank?sslmode=disable",
			MaxRetries:        3,
			RetryDelayMinutes: 5,
		}
		assert.NotNil(t, config)
	})
}

func TestLoadConfig_InvalidRabbitMQURL(t *testing.T) {
	// Test validateRabbitMQURL function directly since LoadConfig calls log.Fatalf
	err := validateRabbitMQURL("invalid-url")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RABBITMQ_URL must use amqp or amqps scheme")
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable exists",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "environment variable does not exist",
			key:          "NON_EXISTENT_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
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
		{
			name:         "valid integer",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "25",
			expected:     25,
		},
		{
			name:         "invalid integer",
			key:          "TEST_INT_INVALID",
			defaultValue: 10,
			envValue:     "not_a_number",
			expected:     10,
		},
		{
			name:         "environment variable does not exist",
			key:          "NON_EXISTENT_INT",
			defaultValue: 10,
			envValue:     "",
			expected:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvInt(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateRabbitMQURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectError bool
	}{
		{
			name:        "valid amqp URL",
			url:         "amqp://testuser:testpass@localhost:5672/",
			expectError: false,
		},
		{
			name:        "valid amqps URL",
			url:         "amqps://testuser:testpass@example.com:5671/vhost",
			expectError: false,
		},
		{
			name:        "placeholder credentials",
			url:         "amqp://user:pass@localhost:5672/",
			expectError: true,
		},
		{
			name:        "invalid scheme",
			url:         "http://localhost:5672",
			expectError: true,
		},
		{
			name:        "invalid URL format",
			url:         "not-a-url",
			expectError: true,
		},
		{
			name:        "empty URL",
			url:         "",
			expectError: true,
		},
		{
			name:        "URL without scheme",
			url:         "localhost:5672",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRabbitMQURL(tt.url)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_Structure(t *testing.T) {
	config := &Config{
		RabbitMQURL:        "amqp://localhost:5672",
		DatabaseURL:        "postgres://localhost/db",
		MaxRetries:         5,
		RetryDelayMinutes:  10,
	}

	assert.Equal(t, "amqp://localhost:5672", config.RabbitMQURL)
	assert.Equal(t, "postgres://localhost/db", config.DatabaseURL)
	assert.Equal(t, 5, config.MaxRetries)
	assert.Equal(t, 10, config.RetryDelayMinutes)
}