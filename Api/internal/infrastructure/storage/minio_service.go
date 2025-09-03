package storage

import (
	"context"
	"io"

	"main_videork/internal/domain/interfaces"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioConfig holds the necessary configuration for connecting to MinIO.
type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

type videoStorage struct {
	client *minio.Client
	bucket string
}

// NewMinioVideoStorage creates a new VideoStorage backed by MinIO.
func NewMinioVideoStorage(cfg MinioConfig) (interfaces.VideoStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &videoStorage{client: client, bucket: cfg.Bucket}, nil
}

// Save uploads the provided video data to MinIO and returns the object name.
func (s *videoStorage) Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	return objectName, nil
}
