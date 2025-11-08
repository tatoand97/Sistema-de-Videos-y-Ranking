package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSPublisher struct {
	client *sqs.Client
}

func NewSQSPublisher(region string) (*SQSPublisher, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return &SQSPublisher{
		client: sqs.NewFromConfig(cfg),
	}, nil
}

// Publish sends the provided body to the queue URL to satisfy interfaces.MessagePublisher.
func (p *SQSPublisher) Publish(queueURL string, body []byte) error {
	if p == nil {
		return fmt.Errorf("sqs publisher is nil")
	}
	if strings.TrimSpace(queueURL) == "" {
		return fmt.Errorf("queue url is required")
	}
	_, err := p.client.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(body)),
	})
	return err
}

// PublishMessage marshals the provided message as JSON and publishes it.
func (p *SQSPublisher) PublishMessage(queueURL string, message interface{}) error {
	messageBody, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return p.Publish(queueURL, messageBody)
}

func (p *SQSPublisher) Close() error {
	return nil
}
