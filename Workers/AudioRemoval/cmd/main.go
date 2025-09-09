package main

import (
	"audioremoval/internal/infrastructure"
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
	defer container.Publisher.Close()

	if err := container.Consumer.StartConsuming(config.QueueName, container.MessageHandler); err != nil {
		logrus.Fatal("Failed to start consuming:", err)
	}

	logrus.Info("AudioRemoval worker started. Waiting for messages...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Consumer.Close()
	logrus.Info("Shutting down AudioRemoval worker...")
}