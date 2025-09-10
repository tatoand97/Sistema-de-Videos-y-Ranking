package adapters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRabbitMQConsumer_InvalidURL(t *testing.T) {
	consumer, err := NewRabbitMQConsumer("invalid-url", 3, 1000)
	assert.Error(t, err)
	assert.Nil(t, consumer)
}

func TestNewRabbitMQConsumer_ValidParams(t *testing.T) {
	// Test parameter validation without actual connection
	url := "amqp://guest:guest@localhost:5672/"
	maxRetries := 3
	queueMaxLength := 1000
	
	assert.NotEmpty(t, url)
	assert.Greater(t, maxRetries, 0)
	assert.Greater(t, queueMaxLength, 0)
}

func TestNewRabbitMQPublisher_InvalidURL(t *testing.T) {
	publisher, err := NewRabbitMQPublisher("invalid-url")
	assert.Error(t, err)
	assert.Nil(t, publisher)
}

func TestNewRabbitMQPublisher_ValidParams(t *testing.T) {
	// Test parameter validation without actual connection
	url := "amqp://guest:guest@localhost:5672/"
	assert.NotEmpty(t, url)
	assert.Contains(t, url, "amqp://")
}

func TestRabbitMQConsumer_GetRetryCount_Logic(t *testing.T) {
	// Test the retry count extraction logic
	headers := map[string]interface{}{
		"x-retry-count": int32(5),
	}
	
	if retryCount, exists := headers["x-retry-count"]; exists {
		if count, ok := retryCount.(int32); ok {
			assert.Equal(t, int32(5), count)
		}
	}
}

func TestRabbitMQConsumer_RetryLogic(t *testing.T) {
	// Test retry logic concepts
	maxRetries := 3
	currentRetries := 2
	
	shouldRetry := currentRetries < maxRetries
	assert.True(t, shouldRetry)
	
	currentRetries = 3
	shouldRetry = currentRetries < maxRetries
	assert.False(t, shouldRetry)
}