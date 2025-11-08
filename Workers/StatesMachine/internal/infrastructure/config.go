package infrastructure

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AWSRegion         string
	SQSQueueURL       string
	TrimVideoQueue    string
	EditVideoQueue    string
	AudioRemovalQueue string
	WatermarkingQueue string
	GossipQueue       string
	DatabaseURL       string
	MaxRetries        int
	RetryDelayMinutes int
	ProcessedVideoURL string
}

func LoadConfig() *Config {
	return &Config{
		AWSRegion:         getEnv("AWS_REGION", "us-east-1"),
		SQSQueueURL:       getEnv("SQS_QUEUE_URL", ""),
		TrimVideoQueue:    getEnv("SQS_TRIM_VIDEO_QUEUE", ""),
		EditVideoQueue:    getEnv("SQS_EDIT_VIDEO_QUEUE", ""),
		AudioRemovalQueue: getEnv("SQS_AUDIO_REMOVAL_QUEUE", ""),
		WatermarkingQueue: getEnv("SQS_WATERMARKING_QUEUE", ""),
		GossipQueue:       getEnv("SQS_GOSSIP_QUEUE", ""),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://app_user:app_password@postgres:5432/videorank?sslmode=disable"),
		MaxRetries:        getEnvInt("MAX_RETRIES", 3),
		RetryDelayMinutes: getEnvInt("RETRY_DELAY_MINUTES", 5),
		ProcessedVideoURL: buildProcessedVideoURL(),
	}
}

func buildProcessedVideoURL() string {
	base := strings.TrimRight(getEnv("PROCESSED_VIDEO_BASE_URL", ""), "/")
	if base != "" {
		return fmt.Sprintf("%s/%%s", base)
	}
	host := getEnv("PROCESSED_VIDEO_HOST", "localhost")
	port := getEnv("PROCESSED_VIDEO_PORT", "8084")
	return fmt.Sprintf("http://%s:%s/processed-videos/%%s", host, port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}


