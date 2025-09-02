package domain

type VideoProcessingService interface {
	InjectCurtains(inputData []byte, curtainInPath, curtainOutPath string) ([]byte, error)
}

type NotificationService interface {
	NotifyProcessingComplete(videoID string, success bool) error
}