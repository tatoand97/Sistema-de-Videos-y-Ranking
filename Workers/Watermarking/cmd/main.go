package main

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"watermarking/internal/infrastructure"
)

func main() {
	_ = godotenv.Load()

	config := infrastructure.LoadConfig()
	container, err := infrastructure.NewContainer(config)
	if err != nil { logrus.Fatal("bootstrap error:", err) }
	defer container.Consumer.Close()

	if err := container.Consumer.StartConsuming(config.QueueName, container.MessageHandler); err != nil {
		logrus.Fatal("start consuming:", err)
	}

	logrus.Info("Watermarking worker started. Waiting for messages...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down Watermarking worker...")
}
