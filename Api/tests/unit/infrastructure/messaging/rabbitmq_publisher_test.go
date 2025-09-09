package messaging_test

import (
	infra "api/internal/infrastructure/messaging"
	"testing"
)

func TestRabbitMQPublisher_Close_NoConn_NoPanic(t *testing.T) {
	var p infra.RabbitMQPublisher
	if err := p.Close(); err != nil {
		t.Fatalf("close returned error: %v", err)
	}
}

func TestRabbitMQPublisher_PublishJSON_MarshalError(t *testing.T) {
	var p infra.RabbitMQPublisher
	// functions are not JSON-marshalable
	err := p.PublishJSON("queue", func() {})
	if err == nil {
		t.Fatalf("expected marshal error, got nil")
	}
}
