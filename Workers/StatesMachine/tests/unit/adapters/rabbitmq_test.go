package adapters

import (
	"statesmachine/internal/adapters"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRabbitMQPublisher_InvalidURL(t *testing.T) {
	publisher, err := adapters.NewRabbitMQPublisher("invalid-url")
	assert.Error(t, err)
	assert.Nil(t, publisher)
}

func TestNewRabbitMQConsumer_InvalidURL(t *testing.T) {
	consumer, err := adapters.NewRabbitMQConsumer("invalid-url")
	assert.Error(t, err)
	assert.Nil(t, consumer)
}