package domain

type VideoProcessingService interface {
	RemoveAudio(inputData []byte) ([]byte, error)
}

type NotificationService interface {
	NotifyVideoProcessed(videoID, filename, bucketPath string) error
	NotifyProcessingComplete(videoID string, success bool) error
}