package main

import (
	"editvideo/internal/infrastructure"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	config := infrastructure.LoadConfig()
	container, err := infrastructure.NewContainer(config)
	if err != nil {
		logrus.Fatal("bootstrap error:", err)
	}
	defer container.Consumer.Close()

	if err := container.Consumer.StartConsuming(func(data []byte) error {
		return container.MessageHandler.HandleMessage(data)
	}); err != nil {
		logrus.Fatal("start consuming:", err)
	}

	logrus.Info("EditVideo worker started. Waiting for messages...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down EditVideo worker...")
}
