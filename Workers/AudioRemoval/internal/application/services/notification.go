package services

import "github.com/sirupsen/logrus"

type LogNotificationService struct{}

func NewLogNotificationService() *LogNotificationService {
	return &LogNotificationService{}
}

func (s *LogNotificationService) NotifyProcessingComplete(videoID string, success bool) error {
	if success {
		logrus.Infof("Video %s processed successfully", videoID)
	} else {
		logrus.Errorf("Video %s processing failed", videoID)
	}
	return nil
}