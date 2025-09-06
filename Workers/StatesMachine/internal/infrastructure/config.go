package infrastructure

import (
	"os"
	"log"
	"fmt"
	"net/url"
)

type Config struct {
	RabbitMQURL         string
	QueueName           string
	EditVideoQueue      string
	AudioRemovalQueue   string
	WatermarkingQueue   string
}

func LoadConfig() *Config {
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://user:pass@rabbitmq:5672/")
	
	if err := validateRabbitMQURL(rabbitURL); err != nil {
		log.Fatalf("RabbitMQ URL validation failed: %v", err)
	}
	
	return &Config{
		RabbitMQURL:       rabbitURL,
		QueueName:         getEnv("QUEUE_NAME", "orders"),
		EditVideoQueue:    getEnv("EDIT_VIDEO_QUEUE", "edit_video_queue"),
		AudioRemovalQueue: getEnv("AUDIO_REMOVAL_QUEUE", "audio_removal_queue"),
		WatermarkingQueue: getEnv("WATERMARKING_QUEUE", "watermarking_queue"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// validateRabbitMQURL validates RabbitMQ connection string
func validateRabbitMQURL(rabbitURL string) error {
	if rabbitURL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}
	
	u, err := url.Parse(rabbitURL)
	if err != nil {
		return fmt.Errorf("invalid RABBITMQ_URL format: %v", err)
	}
	
	if u.Scheme != "amqp" && u.Scheme != "amqps" {
		return fmt.Errorf("RABBITMQ_URL must use amqp or amqps scheme")
	}
	
	// Check for placeholder credentials
	if u.User != nil {
		username := u.User.Username()
		password, _ := u.User.Password()
		if username == "user" && password == "pass" {
			return fmt.Errorf("RABBITMQ_URL contains placeholder credentials, please set real credentials")
		}
	}
	
	return nil
}