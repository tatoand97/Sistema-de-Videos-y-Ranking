package infrastructure

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	RabbitMQURL       string
	S3Region          string
	S3Endpoint        string
	S3AccessKey       string
	S3SecretKey       string
	S3UsePathStyle    bool
	RawBucket         string
	ProcessedBucket   string
	QueueName         string
	StateMachineQueue string
	MaxRetries        int
	QueueMaxLength    int
	MaxSeconds        int
}

func LoadConfig() *Config {
	maxRetries := 3
	if v := os.Getenv("MAX_RETRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			maxRetries = n
		}
	}
	queueMax := 1000
	if v := os.Getenv("QUEUE_MAX_LENGTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			queueMax = n
		}
	}
	maxSeconds := 30
	if v := os.Getenv("MAX_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			maxSeconds = n
		}
	}
	usePathStyle := strings.EqualFold(os.Getenv("S3_USE_PATH_STYLE"), "true")

	return &Config{
		RabbitMQURL:       os.Getenv("RABBITMQ_URL"),
		S3Region:          os.Getenv("AWS_REGION"),
		S3Endpoint:        os.Getenv("S3_ENDPOINT"),
		S3AccessKey:       os.Getenv("AWS_ACCESS_KEY_ID"),
		S3SecretKey:       os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3UsePathStyle:    usePathStyle,
		RawBucket:         os.Getenv("RAW_BUCKET"),
		ProcessedBucket:   os.Getenv("PROCESSED_BUCKET"),
		QueueName:         os.Getenv("QUEUE_NAME"),
		StateMachineQueue: os.Getenv("STATE_MACHINE_QUEUE"),
		MaxRetries:        maxRetries,
		QueueMaxLength:    queueMax,
		MaxSeconds:        maxSeconds,
	}
}
