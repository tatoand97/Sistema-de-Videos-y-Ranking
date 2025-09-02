package infrastructure

import (
	"os"
	"strconv"
)

type Config struct {
	RabbitMQURL     string
	MinIOEndpoint   string
	MinIOAccessKey  string
	MinIOSecretKey  string
	RawBucket       string
	ProcessedBucket string
	QueueName       string
	MaxRetries      int
	QueueMaxLength  int
	MaxSeconds      int
}

func LoadConfig() *Config {
	maxRetries := 3
	if v := os.Getenv("MAX_RETRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil { maxRetries = n }
	}
	queueMax := 1000
	if v := os.Getenv("QUEUE_MAX_LENGTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil { queueMax = n }
	}
	maxSeconds := 30
	if v := os.Getenv("MAX_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil { maxSeconds = n }
	}
	return &Config{
		RabbitMQURL:     os.Getenv("RABBITMQ_URL"),
		MinIOEndpoint:   os.Getenv("MINIO_ENDPOINT"),
		MinIOAccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
		MinIOSecretKey:  os.Getenv("MINIO_SECRET_KEY"),
		RawBucket:       os.Getenv("RAW_BUCKET"),
		ProcessedBucket: os.Getenv("PROCESSED_BUCKET"),
		QueueName:       os.Getenv("QUEUE_NAME"),
		MaxRetries:      maxRetries,
		QueueMaxLength:  queueMax,
		MaxSeconds:      maxSeconds,
	}
}
