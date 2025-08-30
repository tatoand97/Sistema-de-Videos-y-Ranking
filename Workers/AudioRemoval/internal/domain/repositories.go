package domain

type VideoRepository interface {
	FindByFilename(filename string) (*Video, error)
	UpdateStatus(id string, status ProcessingStatus) error
}

type StorageRepository interface {
	Download(bucket, filename string) ([]byte, error)
	Upload(bucket, filename string, data []byte) error
}