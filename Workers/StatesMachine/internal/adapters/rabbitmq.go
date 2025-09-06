package adapters

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/sirupsen/logrus"
)

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQConsumer(url string) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil { return nil, err }

	ch, err := conn.Channel()
	if err != nil { return nil, err }

	return &RabbitMQConsumer{conn: conn, channel: ch}, nil
}

func NewRabbitMQPublisher(url string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil { return nil, err }

	ch, err := conn.Channel()
	if err != nil { return nil, err }

	return &RabbitMQPublisher{conn: conn, channel: ch}, nil
}

func (r *RabbitMQConsumer) StartConsuming(queueName string, handler MessageHandlerInterface) error {
	q, err := r.channel.QueueDeclare(queueName, true, false, false, false, amqp.Table{
		"x-max-length": 1000,
	})
	if err != nil { return err }

	msgs, err := r.channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil { return err }

	logrus.Infof("Consumiendo cola %s (maxlen=1000)", queueName)

	go func() {
		for d := range msgs {
			if err := handler.HandleMessage(d.Body); err != nil {
				logrus.Errorf("Error processing message: %v", err)
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	return nil
}

func (r *RabbitMQPublisher) PublishMessage(queueName string, message []byte) error {
	for attempts := 0; attempts < 3; attempts++ {
		if r.channel == nil {
			if err := r.reconnect(); err != nil {
				logrus.Errorf("Reconnect attempt %d failed: %v", attempts+1, err)
				continue
			}
		}

		// Simple queue configuration for all queues
		args := amqp.Table{
			"x-max-length": 1000,
		}

		_, err := r.channel.QueueDeclare(queueName, true, false, false, false, args)
		if err != nil {
			logrus.Errorf("Queue declare failed: %v", err)
			r.channel = nil
			continue
		}

		err = r.channel.Publish("", queueName, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
		if err != nil {
			logrus.Errorf("Publish failed: %v", err)
			r.channel = nil
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to publish after 3 attempts")
}

func (r *RabbitMQPublisher) reconnect() error {
	if r.conn != nil {
		r.conn.Close()
	}

	conn, err := amqp.Dial("amqp://admin:admin@rabbitmq:5672/")
	if err != nil { return err }

	ch, err := conn.Channel()
	if err != nil { return err }

	r.conn = conn
	r.channel = ch
	return nil
}





func (r *RabbitMQConsumer) Close() error {
	if r.channel != nil { r.channel.Close() }
	if r.conn != nil { r.conn.Close() }
	return nil
}

type MessageHandlerInterface interface {
	HandleMessage(body []byte) error
}