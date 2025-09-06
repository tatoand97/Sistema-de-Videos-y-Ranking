package interfaces

// MessagePublisher abstracts a message broker publisher.
// Implemented in infra (e.g., RabbitMQPublisher).
type MessagePublisher interface {
	Publish(queue string, body []byte) error
	Close() error
}
