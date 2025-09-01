package services

import "github.com/sirupsen/logrus"

type LogNotificationService struct{}

func NewLogNotificationService() *LogNotificationService { return &LogNotificationService{} }

func (s *LogNotificationService) NotifyProcessingComplete(videoID string, success bool) error {
    if success {
        logrus.Infof("[notify] video %s procesado con éxito", videoID)
    } else {
        logrus.Warnf("[notify] video %s falló", videoID)
    }
    return nil
}