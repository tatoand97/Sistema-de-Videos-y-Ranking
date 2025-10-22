package infrastructure_test

import (
	"audioremoval/internal/infrastructure"
	"testing"
)

func TestNewContainer_InvalidS3Config(t *testing.T) {
	cfg := &infrastructure.Config{
		RabbitMQURL:     "amqp://localhost:5672",
		RawBucket:       "raw",
		ProcessedBucket: "processed",
	}

	container, err := infrastructure.NewContainer(cfg)
	if err == nil || container != nil {
		t.Fatalf("expected error due to missing S3 region, got container=%v err=%v", container, err)
	}
}
