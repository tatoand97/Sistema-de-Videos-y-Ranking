package messaging

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSConsumer struct {
	client   *sqs.Client
	queueURL string
}

type MessageHandler func([]byte) error

func NewSQSConsumer(region, queueURL string) (*SQSConsumer, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(os.Getenv("AWS_PROFILE")),
	)
	if err != nil {
		return nil, err
	}

	return &SQSConsumer{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}

func (c *SQSConsumer) StartConsuming(handler MessageHandler) error {
	for {
		result, err := c.client.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(c.queueURL),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     20,
		})
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, message := range result.Messages {
			if err := handler([]byte(*message.Body)); err != nil {
				log.Printf("Error processing message: %v", err)
				continue
			}

			_, err := c.client.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(c.queueURL),
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Printf("Error deleting message: %v", err)
			}
		}
	}
}

func (c *SQSConsumer) PublishMessage(queueURL string, message interface{}) error {
	messageBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = c.client.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(messageBody)),
	})
	return err
}

// Close is a no-op placeholder to satisfy publisher interfaces.
func (c *SQSConsumer) Close() error {
	return nil
}
