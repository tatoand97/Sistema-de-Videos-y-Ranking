package domain

type VideoRepository interface {
	FindByID(id uint) (*Video, error)
	UpdateStatus(id uint, status VideoStatus) error
	UpdateStatusAndProcessedFile(id uint, status VideoStatus, processedFile string) error
}

type MessagePublisher interface {
	PublishMessage(queueName string, message []byte) error
}