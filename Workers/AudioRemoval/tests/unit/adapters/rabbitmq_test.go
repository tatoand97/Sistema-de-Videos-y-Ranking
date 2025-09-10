package adapters_test

import (
	"audioremoval/internal/adapters"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRabbitMQConsumer(t *testing.T) {
	consumer := adapters.NewRabbitMQConsumer("amqp://localhost", "test-queue")
	assert.NotNil(t, consumer)
}

func TestNewRabbitMQPublisher(t *testing.T) {
	publisher := adapters.NewRabbitMQPublisher("amqp://localhost")
	assert.NotNil(t, publisher)
}

func TestRabbitMQConsumer_StartConsuming_InvalidURL(t *testing.T) {
	consumer := adapters.NewRabbitMQConsumer("invalid-url", "test-queue")
	
	err := consumer.StartConsuming(nil)
	
	assert.Error(t, err)
}

func TestRabbitMQPublisher_PublishMessage_InvalidURL(t *testing.T) {
	publisher := adapters.NewRabbitMQPublisher("invalid-url")
	
	err := publisher.PublishMessage("test-queue", "test message")
	
	assert.Error(t, err)
}

func TestRabbitMQConsumer_Close(t *testing.T) {
	consumer := adapters.NewRabbitMQConsumer("amqp://localhost", "test-queue")
	
	err := consumer.Close()
	
	// Should not error even if not connected
	assert.NoError(t, err)
}

func TestRabbitMQPublisher_Close(t *testing.T) {
	publisher := adapters.NewRabbitMQPublisher("amqp://localhost")
	
	err := publisher.Close()
	
	// Should not error even if not connected
	assert.NoError(t, err)
}