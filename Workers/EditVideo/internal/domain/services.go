package domain

type VideoProcessingService interface {
	TrimToMaxSeconds(inputData []byte, maxSeconds int) ([]byte, error)
}

type NotificationService interface {
	NotifyProcessingComplete(videoID string, success bool) error
}
