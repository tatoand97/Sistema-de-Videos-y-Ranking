package adapters

import (
	"context"
	"io"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStorage struct {
	client *minio.Client
}

func NewMinIOStorage(endpoint, accessKey, secretKey string) (*MinIOStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil { return nil, err }
	return &MinIOStorage{client: client}, nil
}

func (m *MinIOStorage) GetObject(bucket, filename string) (io.Reader, error) {
	obj, err := m.client.GetObject(context.Background(), bucket, filename, minio.GetObjectOptions{})
	if err != nil { return nil, err }
	return obj, nil
}

func (m *MinIOStorage) PutObject(bucket, filename string, data io.Reader, size int64) error {
	_, err := m.client.PutObject(context.Background(), bucket, filename, data, size, minio.PutObjectOptions{})
	return err
}
