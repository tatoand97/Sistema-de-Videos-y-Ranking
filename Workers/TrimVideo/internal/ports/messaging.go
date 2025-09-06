package ports

type MessageHandler interface {
	HandleMessage(body []byte) error
}

type MessageConsumer interface {
	StartConsuming(queueName string, handler MessageHandler) error
	Close() error
}

type MessagePublisher interface {
	PublishMessage(queueName string, message interface{}) error
	Close() error
}
