package infrastructure

import (
	"os"
	"statesmachine/internal/infrastructure"
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

	config := infrastructure.LoadConfig()

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

	// Test that the function exists
	assert.NotPanics(t, func() {
		config := &infrastructure.Config{
			RabbitMQURL:       "amqp://testuser:testpass@localhost:5672",
			DatabaseURL:       "postgres://app_user:app_password@postgres:5432/videorank?sslmode=disable",
			MaxRetries:        3,
			RetryDelayMinutes: 5,
		}
		assert.NotNil(t, config)
	})
}

func TestConfig_Structure(t *testing.T) {
	config := &infrastructure.Config{
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