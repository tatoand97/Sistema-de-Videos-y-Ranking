package main

import (
	"curtaininjector/internal/adapters"
	"curtaininjector/internal/infrastructure"
	"os"
	"os/signal"
	"syscall"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()

	config := infrastructure.LoadConfig()
	container, err := infrastructure.NewContainer(config)
	if err != nil {
		logrus.Fatal("Failed to initialize container:", err)
	}

	messageHandler := adapters.NewMessageHandler(container.ProcessVideoUC)
	consumer, err := adapters.NewRabbitMQConsumer(config.RabbitMQURL, config.MaxRetries, config.QueueMaxLength)
	if err != nil {
		logrus.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer consumer.Close()

	if err := consumer.StartConsuming(config.QueueName, messageHandler); err != nil {
		logrus.Fatal("Failed to start consuming:", err)
	}

	logrus.Info("CurtainInjector worker started. Waiting for messages...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down CurtainInjector worker...")
}