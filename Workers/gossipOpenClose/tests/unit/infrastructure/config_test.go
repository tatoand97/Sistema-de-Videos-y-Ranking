package infrastructure

import (
	"gossipopenclose/internal/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Success(t *testing.T) {
	config := infrastructure.LoadConfig()
	assert.NotNil(t, config)
}

func TestConfig_Structure(t *testing.T) {
	config := &infrastructure.Config{
		RabbitMQURL: "amqp://localhost:5672",
		MinIOEndpoint: "localhost:9000",
	}
	
	assert.Equal(t, "amqp://localhost:5672", config.RabbitMQURL)
	assert.Equal(t, "localhost:9000", config.MinIOEndpoint)
}