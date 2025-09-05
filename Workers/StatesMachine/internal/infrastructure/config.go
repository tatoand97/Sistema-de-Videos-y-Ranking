package infrastructure

import "os"

type Config struct {
	RabbitMQURL       string
	QueueName         string
	EditVideoQueue    string
	AudioRemovalQueue string
}

func LoadConfig() *Config {
	return &Config{
		RabbitMQURL:       getEnv("RABBITMQ_URL", "amqp://admin:admin@rabbitmq:5672/"),
		QueueName:         getEnv("QUEUE_NAME", "orders"),
		EditVideoQueue:    getEnv("EDIT_VIDEO_QUEUE", "edit_video_queue"),
		AudioRemovalQueue: getEnv("AUDIO_REMOVAL_QUEUE", "audio_removal_queue"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}