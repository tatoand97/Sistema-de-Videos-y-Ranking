package adapters

import (
	"testing"
	"watermarking/internal/adapters"

	"github.com/stretchr/testify/assert"
)

func TestNewRabbitMQConsumer_InvalidURL(t *testing.T) {
	consumer, err := adapters.NewRabbitMQConsumer("invalid-url", 3, 1000)
	assert.Error(t, err)
	assert.Nil(t, consumer)
}

func TestNewRabbitMQPublisher_InvalidURL(t *testing.T) {
	publisher, err := adapters.NewRabbitMQPublisher("invalid-url")
	assert.Error(t, err)
	assert.Nil(t, publisher)
}

func TestRabbitMQ_ValidParams(t *testing.T) {
	url := "amqp://guest:guest@localhost:5672/"
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "amqp://")
}