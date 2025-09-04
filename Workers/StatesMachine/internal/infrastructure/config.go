package infrastructure

import "os"

type Config struct {
	RabbitMQURL string
	QueueName   string
}

func LoadConfig() *Config {
	return &Config{
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://admin:admin@rabbitmq:5672/"),
		QueueName:   getEnv("QUEUE_NAME", "states_machine_queue"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}