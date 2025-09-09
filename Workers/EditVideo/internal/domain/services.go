package domain

type VideoProcessingService interface {
	TrimToMaxSeconds(inputData []byte, maxSeconds int) ([]byte, error)
}

type NotificationService interface {
	NotifyVideoProcessed(videoID, filename, bucketPath string) error
}
