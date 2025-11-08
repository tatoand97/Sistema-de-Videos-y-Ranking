package messaging

import (
	"context"
	"encoding/json"

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

func (p *SQSPublisher) PublishMessage(queueURL string, message interface{}) error {
	messageBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = p.client.SendMessage(context.Background(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(messageBody)),
	})
	return err
}

func (p *SQSPublisher) Close() error {
	return nil
}