// internal/infrastructure/config.go
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
	IntroSeconds    float64
	OutroSeconds    float64
	TargetWidth     int
	TargetHeight    int
	FPS             int
	LogoPath        string
}

func LoadConfig() *Config {
	maxRetries := getEnvInt("MAX_RETRIES", 5)
	queueMax := getEnvInt("QUEUE_MAX_LENGTH", 1000)
	maxSeconds := getEnvInt("MAX_SECONDS", 30)

	intro := getEnvFloat("INTRO_SECONDS", 2.5)
	outro := getEnvFloat("OUTRO_SECONDS", 2.5)
	tw := getEnvInt("TARGET_WIDTH", 1280)
	th := getEnvInt("TARGET_HEIGHT", 720)
	fps := getEnvInt("FPS", 30)

	logo := os.Getenv("LOGO_PATH")
	if logo == "" {
		logo = "./assets/nba-logo-removebg-preview.png"
	}

	return &Config{
		// Usa RABBITMQ_URL
		RabbitMQURL:     os.Getenv("RABBITMQ_URL"),
		MinIOEndpoint:   os.Getenv("MINIO_ENDPOINT"),
		MinIOAccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
		MinIOSecretKey:  os.Getenv("MINIO_SECRET_KEY"),
		RawBucket:       os.Getenv("MINIO_BUCKET_RAW"),
		ProcessedBucket: os.Getenv("MINIO_BUCKET_PROCESSED"),
		QueueName:       os.Getenv("QUEUE_NAME"),
		MaxRetries:      maxRetries,
		QueueMaxLength:  queueMax,
		MaxSeconds:      maxSeconds,
		IntroSeconds:    intro,
		OutroSeconds:    outro,
		TargetWidth:     tw,
		TargetHeight:    th,
		FPS:             fps,
		LogoPath:        logo,
	}
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getEnvFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}