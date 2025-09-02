package services

import (
	"github.com/sirupsen/logrus"
)

type LogNotificationService struct{}

func NewLogNotificationService() *LogNotificationService {
	return &LogNotificationService{}
}

func (s *LogNotificationService) NotifyProcessingComplete(videoID string, success bool) error {
	if success {
		logrus.WithField("videoID", videoID).Info("Video curtain injection completed successfully")
	} else {
		logrus.WithField("videoID", videoID).Error("Video curtain injection failed")
	}
	return nil
}