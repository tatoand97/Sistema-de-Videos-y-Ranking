package domain

type VideoRepository interface {
	FindByID(id uint) (*Video, error)
	UpdateStatus(id uint, status VideoStatus) error
}

type MessagePublisher interface {
	PublishMessage(queueName string, message []byte) error
}