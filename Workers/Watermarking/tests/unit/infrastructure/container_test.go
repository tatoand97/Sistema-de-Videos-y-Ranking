package infrastructure

import (
	"testing"
	"watermarking/internal/infrastructure"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer_NilConfig(t *testing.T) {
	assert.Panics(t, func() {
		infrastructure.NewContainer(nil)
	})
}

func TestNewContainer_InvalidMinIOConfig(t *testing.T) {
	config := &infrastructure.Config{
		S3Region:    "",
		RabbitMQURL: "amqp://localhost:5672",
	}

	container, err := infrastructure.NewContainer(config)
	assert.Error(t, err)
	assert.Nil(t, container)
}

func TestNewContainer_InvalidRabbitMQConfig(t *testing.T) {
	config := &infrastructure.Config{
		S3Region:    "us-east-1",
		RabbitMQURL: "invalid-url",
	}

	container, err := infrastructure.NewContainer(config)
	assert.Error(t, err)
	assert.Nil(t, container)
}

func TestContainer_Structure(t *testing.T) {
	container := &infrastructure.Container{
		Config:         &infrastructure.Config{},
		Consumer:       nil,
		MessageHandler: nil,
	}

	assert.NotNil(t, container)
	assert.NotNil(t, container.Config)
}
