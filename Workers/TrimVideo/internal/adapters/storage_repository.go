package adapters

import (
	"trimvideo/internal/ports"
	"bytes"
	"io"
)

type StorageRepository struct {
	storage ports.StorageService
}

func NewStorageRepository(storage ports.StorageService) *StorageRepository {
	return &StorageRepository{storage: storage}
}

func (r *StorageRepository) Download(bucket, filename string) ([]byte, error) {
	reader, err := r.storage.GetObject(bucket, filename)
	if err != nil { return nil, err }
	return io.ReadAll(reader)
}

func (r *StorageRepository) Upload(bucket, filename string, data []byte) error {
	return r.storage.PutObject(bucket, filename, bytes.NewReader(data), int64(len(data)))
}
