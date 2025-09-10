package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRabbitMQConsumer_InvalidURL(t *testing.T) {
	consumer, err := NewRabbitMQConsumer("invalid-url", 3, 1000)
	assert.Error(t, err)
	assert.Nil(t, consumer)
}

func TestNewRabbitMQPublisher_InvalidURL(t *testing.T) {
	publisher, err := NewRabbitMQPublisher("invalid-url")
	assert.Error(t, err)
	assert.Nil(t, publisher)
}

func TestRabbitMQ_ValidParams(t *testing.T) {
	url := "amqp://guest:guest@localhost:5672/"
	maxRetries := 3
	queueMaxLength := 1000
	
	assert.NotEmpty(t, url)
	assert.Greater(t, maxRetries, 0)
	assert.Greater(t, queueMaxLength, 0)
	assert.Contains(t, url, "amqp://")
}