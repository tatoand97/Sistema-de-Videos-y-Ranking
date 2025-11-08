package infrastructure

import (
	"encoding/json"

	"shared/messaging"
)

// SQSPublisherAdapter adapts the shared SQS consumer to the MessagePublisher interface.
type SQSPublisherAdapter struct {
	client *messaging.SQSConsumer
}

func NewSQSPublisherAdapter(client *messaging.SQSConsumer) *SQSPublisherAdapter {
	return &SQSPublisherAdapter{client: client}
}

func (a *SQSPublisherAdapter) PublishMessage(queueName string, message []byte) error {
	raw := json.RawMessage(message)
	return a.client.PublishMessage(queueName, raw)
}
