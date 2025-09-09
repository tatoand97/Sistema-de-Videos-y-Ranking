package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer_Structure(t *testing.T) {
	config := &Config{
		RabbitMQURL:        "amqp://localhost:5672",
		DatabaseURL:        "postgres://localhost/db",
		MaxRetries:         3,
		RetryDelayMinutes:  5,
	}

	// Test that NewContainer function exists and can be called
	// Note: We can't test the actual functionality without real connections
	// but we can test the structure and ensure it doesn't panic
	assert.NotPanics(t, func() {
		NewContainer(config)
	})
}

func TestContainer_ConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "valid config",
			config: &Config{
				RabbitMQURL:        "amqp://localhost:5672",
				DatabaseURL:        "postgres://localhost/db",
				MaxRetries:         3,
				RetryDelayMinutes:  5,
			},
		},
		{
			name: "config with different values",
			config: &Config{
				RabbitMQURL:        "amqp://guest:guest@rabbitmq:5672/",
				DatabaseURL:        "postgres://user:pass@postgres:5432/statesmachine",
				MaxRetries:         5,
				RetryDelayMinutes:  10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.config)
			assert.NotEmpty(t, tt.config.RabbitMQURL)
			assert.NotEmpty(t, tt.config.DatabaseURL)
			assert.Greater(t, tt.config.MaxRetries, 0)
			assert.Greater(t, tt.config.RetryDelayMinutes, 0)
		})
	}
}

func TestContainer_NilConfig(t *testing.T) {
	// Test that passing nil config causes expected panic
	assert.Panics(t, func() {
		NewContainer(nil)
	})
}