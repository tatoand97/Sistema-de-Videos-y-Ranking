package infrastructure

import "os"

type Config struct {
	RabbitMQURL     string
	MinIOEndpoint   string
	MinIOAccessKey  string
	MinIOSecretKey  string
	RawBucket       string
	ProcessedBucket string
	QueueName       string
}

func LoadConfig() *Config {
	return &Config{
		RabbitMQURL:     os.Getenv("RABBITMQ_URL"),
		MinIOEndpoint:   os.Getenv("MINIO_ENDPOINT"),
		MinIOAccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
		MinIOSecretKey:  os.Getenv("MINIO_SECRET_KEY"),
		RawBucket:       os.Getenv("RAW_BUCKET"),
		ProcessedBucket: os.Getenv("PROCESSED_BUCKET"),
		QueueName:       os.Getenv("QUEUE_NAME"),
	}
}