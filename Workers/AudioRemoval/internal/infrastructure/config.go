package infrastructure

import (
	"os"
	"strconv"
)

type Config struct {
	RabbitMQURL       string
	MinIOEndpoint     string
	MinIOAccessKey    string
	MinIOSecretKey    string
	RawBucket         string
	ProcessedBucket   string
	QueueName         string
	StateMachineQueue string
	MaxRetries        int
	QueueMaxLength    int
}

func LoadConfig() *Config {
	maxRetries := 3 // default value
	if retries := os.Getenv("MAX_RETRIES"); retries != "" {
		if parsed, err := strconv.Atoi(retries); err == nil {
			maxRetries = parsed
		}
	}
	
	queueMaxLength := 1000 // default value
	if maxLength := os.Getenv("QUEUE_MAX_LENGTH"); maxLength != "" {
		if parsed, err := strconv.Atoi(maxLength); err == nil {
			queueMaxLength = parsed
		}
	}
	
	return &Config{
		RabbitMQURL:       os.Getenv("RABBITMQ_URL"),
		MinIOEndpoint:     os.Getenv("MINIO_ENDPOINT"),
		MinIOAccessKey:    os.Getenv("MINIO_ACCESS_KEY"),
		MinIOSecretKey:    os.Getenv("MINIO_SECRET_KEY"),
		RawBucket:         os.Getenv("RAW_BUCKET"),
		ProcessedBucket:   os.Getenv("PROCESSED_BUCKET"),
		QueueName:         os.Getenv("QUEUE_NAME"),
		StateMachineQueue: os.Getenv("STATE_MACHINE_QUEUE"),
		MaxRetries:        maxRetries,
		QueueMaxLength:    queueMaxLength,
	}
}