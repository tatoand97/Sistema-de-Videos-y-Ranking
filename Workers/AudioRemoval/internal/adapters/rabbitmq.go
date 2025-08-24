package adapters

import (
	"audioremoval/internal/ports"
	"strconv"
	"github.com/streadway/amqp"
	"github.com/sirupsen/logrus"
)

type RabbitMQConsumer struct {
	conn           *amqp.Connection
	channel        *amqp.Channel
	maxRetries     int
	queueMaxLength int
}

func NewRabbitMQConsumer(url string, maxRetries, queueMaxLength int) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:           conn,
		channel:        ch,
		maxRetries:     maxRetries,
		queueMaxLength: queueMaxLength,
	}, nil
}

func (r *RabbitMQConsumer) StartConsuming(queueName string, handler ports.MessageHandler) error {
	// Declare dead letter exchange and queue
	dlxName := queueName + ".dlx"
	dlqName := queueName + ".dlq"
	
	err := r.channel.ExchangeDeclare(dlxName, "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}
	
	_, err = r.channel.QueueDeclare(dlqName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	
	err = r.channel.QueueBind(dlqName, queueName, dlxName, false, nil)
	if err != nil {
		return err
	}

	// Declare main queue with DLX configuration and length limit
	args := amqp.Table{
		"x-dead-letter-exchange": dlxName,
		"x-dead-letter-routing-key": queueName,
		"x-max-length": r.queueMaxLength,
		"x-overflow": "reject-publish-dlx", // Send overflow to DLX
	}
	_, err = r.channel.QueueDeclare(queueName, true, false, false, false, args)
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
				
				retryCount := r.getRetryCount(msg)
				if retryCount >= r.maxRetries {
					logrus.Errorf("Message exceeded max retries (%d), sending to DLQ", r.maxRetries)
					msg.Nack(false, false) // Send to DLQ
				} else {
					logrus.Infof("Retrying message (attempt %d/%d)", retryCount+1, r.maxRetries)
					msg.Nack(false, true) // Requeue
				}
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
func (r *RabbitMQConsumer) getRetryCount(msg amqp.Delivery) int {
	if msg.Headers == nil {
		return 0
	}
	
	if retryCount, exists := msg.Headers["x-retry-count"]; exists {
		if count, ok := retryCount.(string); ok {
			if parsed, err := strconv.Atoi(count); err == nil {
				return parsed
			}
		}
	}
	
	// Check x-death header for retry count from RabbitMQ
	if xDeath, exists := msg.Headers["x-death"]; exists {
		if deaths, ok := xDeath.([]interface{}); ok && len(deaths) > 0 {
			if death, ok := deaths[0].(amqp.Table); ok {
				if count, exists := death["count"]; exists {
					if countVal, ok := count.(int64); ok {
						return int(countVal)
					}
				}
			}
		}
	}
	
	return 0
}