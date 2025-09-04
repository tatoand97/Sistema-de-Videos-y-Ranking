package adapters

import (
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
	q, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
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
	_, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil { return err }

	return r.channel.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
}

func (r *RabbitMQConsumer) PublishMessage(queueName string, message []byte) error {
	return r.channel.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
}

func (r *RabbitMQConsumer) Close() error {
	if r.channel != nil { r.channel.Close() }
	if r.conn != nil { r.conn.Close() }
	return nil
}

type MessageHandlerInterface interface {
	HandleMessage(body []byte) error
}