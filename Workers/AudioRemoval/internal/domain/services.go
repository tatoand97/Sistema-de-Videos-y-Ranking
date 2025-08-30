package domain

type VideoProcessingService interface {
	RemoveAudio(inputData []byte) ([]byte, error)
}

type NotificationService interface {
	NotifyProcessingComplete(videoID string, success bool) error
}