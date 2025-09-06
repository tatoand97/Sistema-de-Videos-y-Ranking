package domain

type VideoRepository interface {
	FindByFilename(filename string) (*Video, error)
	UpdateStatus(id string, status VideoStatus) error
}

type MessagePublisher interface {
	PublishMessage(queueName string, message []byte) error
}