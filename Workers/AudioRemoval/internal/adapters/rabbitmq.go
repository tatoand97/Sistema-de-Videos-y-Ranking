package adapters

import (
	"audioremoval/internal/ports"
	"github.com/streadway/amqp"
	"github.com/sirupsen/logrus"
)

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQConsumer(url string) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQConsumer) StartConsuming(queueName string, handler ports.MessageHandler) error {
	_, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handler.HandleMessage(msg.Body); err != nil {
				logrus.Error("Error processing message:", err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (r *RabbitMQConsumer) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	return nil
}