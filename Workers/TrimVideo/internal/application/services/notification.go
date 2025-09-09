package services

import (
	"github.com/sirupsen/logrus"
	"trimvideo/internal/ports"
)

type VideoProcessedMessage struct {
	VideoID    string `json:"video_id"`
	Filename   string `json:"filename"`
	BucketPath string `json:"bucket_path"`
	Status     string `json:"status"`
}

type NotificationService struct {
	publisher     ports.MessagePublisher
	stateQueue    string
}

func NewNotificationService(publisher ports.MessagePublisher, stateQueue string) *NotificationService {
	return &NotificationService{
		publisher:  publisher,
		stateQueue: stateQueue,
	}
}

func (s *NotificationService) NotifyVideoProcessed(videoID, filename, bucketPath string) error {
	msg := VideoProcessedMessage{
		VideoID:    videoID,
		Filename:   filename,
		BucketPath: bucketPath,
		Status:     "completed",
	}

	if err := s.publisher.PublishMessage(s.stateQueue, msg); err != nil {
		logrus.Errorf("Failed to notify state machine: %v", err)
		return err
	}

	logrus.WithFields(logrus.Fields{
		"video_id": videoID,
		"filename": filename,
		"bucket_path": bucketPath,
	}).Info("Video processing notification sent to state machine")

	return nil
}