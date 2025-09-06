package security

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// ValidateRabbitMQURL validates RabbitMQ connection string
func ValidateRabbitMQURL(rabbitURL string) error {
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

// ValidateMinIOConfig validates MinIO configuration
func ValidateMinIOConfig(endpoint, accessKey, secretKey string) error {
	if endpoint == "" {
		return fmt.Errorf("MINIO_ENDPOINT is required")
	}
	if accessKey == "" {
		return fmt.Errorf("MINIO_ACCESS_KEY is required")
	}
	if secretKey == "" {
		return fmt.Errorf("MINIO_SECRET_KEY is required")
	}
	if len(secretKey) < 8 {
		return fmt.Errorf("MINIO_SECRET_KEY must be at least 8 characters")
	}
	return nil
}

// GetEnvWithValidation gets environment variable with validation
func GetEnvWithValidation(key, defaultValue string, required bool) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		if required {
			return "", fmt.Errorf("environment variable %s is required", key)
		}
		return defaultValue, nil
	}
	return value, nil
}

// GetIntEnvWithValidation gets integer environment variable with validation
func GetIntEnvWithValidation(key string, defaultValue, min, max int) (int, error) {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue, nil
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue, fmt.Errorf("invalid integer value for %s: %v", key, err)
	}
	
	if value < min || value > max {
		return defaultValue, fmt.Errorf("%s must be between %d and %d, got %d", key, min, max, value)
	}
	
	return value, nil
}

// ValidateQueueName validates queue name format
func ValidateQueueName(queueName string) error {
	if queueName == "" {
		return fmt.Errorf("queue name cannot be empty")
	}
	if len(queueName) > 255 {
		return fmt.Errorf("queue name too long (max 255 characters)")
	}
	if strings.Contains(queueName, " ") {
		return fmt.Errorf("queue name cannot contain spaces")
	}
	return nil
}