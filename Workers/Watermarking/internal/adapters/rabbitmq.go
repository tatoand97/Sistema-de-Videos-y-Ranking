package adapters

import (
    "github.com/sirupsen/logrus"
    "github.com/streadway/amqp"
    "watermarking/internal/ports"
)

type RabbitMQConsumer struct {
	conn           *amqp.Connection
	channel        *amqp.Channel
	maxRetries     int
	queueMaxLength int
}

func NewRabbitMQConsumer(url string, maxRetries, queueMaxLength int) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil { return nil, err }
	ch, err := conn.Channel()
	if err != nil { return nil, err }
	return &RabbitMQConsumer{conn: conn, channel: ch, maxRetries: maxRetries, queueMaxLength: queueMaxLength}, nil
}

func (r *RabbitMQConsumer) StartConsuming(queueName string, handler ports.MessageHandler) error {
	// Declaraciones básicas
	args := amqp.Table{"x-max-length": int32(r.queueMaxLength)}
	_, err := r.channel.QueueDeclare(queueName, true, false, false, false, args)
	if err != nil { return err }

	deliveries, err := r.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil { return err }

	logrus.Infof("Consumiendo cola %s (maxlen=%d)", queueName, r.queueMaxLength)

	for msg := range deliveries {
		if err := handler.HandleMessage(msg.Body); err != nil {
			logrus.Errorf("handler error: %v", err)
			_ = msg.Nack(false, false)
			continue
		}
		_ = msg.Ack(false)
	}
	return nil
}

func (r *RabbitMQConsumer) Close() error { 
	if r.channel != nil { _ = r.channel.Close() }
	if r.conn != nil { return r.conn.Close() }
	return nil
}
