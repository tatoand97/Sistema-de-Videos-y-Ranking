package ports

import "io"

type StorageService interface {
	GetObject(bucket, filename string) (io.Reader, error)
	PutObject(bucket, filename string, data io.Reader, size int64) error
}
