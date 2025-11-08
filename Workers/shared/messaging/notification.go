package messaging

type NotificationService struct {
	publisher *SQSConsumer
	queueURL  string
}

func NewNotificationService(publisher *SQSConsumer, queueURL string) *NotificationService {
	return &NotificationService{
		publisher: publisher,
		queueURL:  queueURL,
	}
}

func (n *NotificationService) NotifyCompletion(videoID, filename string) error {
	message := map[string]interface{}{
		"video_id": videoID,
		"filename": filename,
		"status":   "completed",
	}

	return n.publisher.PublishMessage(n.queueURL, message)
}

func (n *NotificationService) NotifyError(videoID, filename, errorMsg string) error {
	message := map[string]interface{}{
		"video_id": videoID,
		"filename": filename,
		"status":   "error",
		"error":    errorMsg,
	}

	return n.publisher.PublishMessage(n.queueURL, message)
}
