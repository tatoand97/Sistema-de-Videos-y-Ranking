package messaging

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQPublisher publishes messages to RabbitMQ queues.
type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQPublisher(url string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	// Publisher confirms can be enabled later if needed
	return &RabbitMQPublisher{conn: conn, channel: ch}, nil
}

// EnsureQueue declares a durable queue with optional DLX and max length.
// It is safe to call multiple times; server will keep existing settings when compatible.
func (p *RabbitMQPublisher) EnsureQueue(queueName string, maxLen int, withDLQ bool) error {
	if withDLQ {
		dlxName := queueName + ".dlx"
		dlqName := queueName + ".dlq"
		if err := p.channel.ExchangeDeclare(dlxName, "direct", true, false, false, false, nil); err != nil {
			return err
		}
		if _, err := p.channel.QueueDeclare(dlqName, true, false, false, false, nil); err != nil {
			return err
		}
		if err := p.channel.QueueBind(dlqName, queueName, dlxName, false, nil); err != nil {
			return err
		}
		args := amqp.Table{
			"x-dead-letter-exchange":    dlxName,
			"x-dead-letter-routing-key": queueName,
			"x-max-length":              maxLen,
			"x-overflow":                "reject-publish-dlx",
		}
		if _, err := p.channel.QueueDeclare(queueName, true, false, false, false, args); err != nil {
			return err
		}
		return nil
	}
	// Basic durable queue without DLX
	_, err := p.channel.QueueDeclare(queueName, true, false, false, false, amqp.Table{"x-max-length": maxLen})
	return err
}

// Publish sends a message with persistent delivery mode and application/json content-type.
func (p *RabbitMQPublisher) Publish(queueName string, body []byte) error {
	return p.channel.Publish("", queueName, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
}

// PublishJSON marshals v to JSON and publishes it.
func (p *RabbitMQPublisher) PublishJSON(queueName string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return p.Publish(queueName, b)
}

func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			log.Printf("rabbitmq publisher channel close: %v", err)
		}
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
