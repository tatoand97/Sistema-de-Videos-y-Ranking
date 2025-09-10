package storage_test

import (
	stor "api/internal/infrastructure/storage"
	"testing"
)

func TestNewMinioVideoStorage_CreatesClient(t *testing.T) {
	cfg := stor.MinioConfig{Endpoint: "play.min.io", AccessKey: "minioadmin", SecretKey: "minioadmin", UseSSL: true, Bucket: "mybucket"}
	s, err := stor.NewMinioVideoStorage(cfg)
	if err != nil || s == nil {
		t.Fatalf("unexpected error creating storage: %v", err)
	}
}
