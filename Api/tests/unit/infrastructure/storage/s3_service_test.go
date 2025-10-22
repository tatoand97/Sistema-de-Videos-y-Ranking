package storage_test

import (
	"testing"

	stor "api/internal/infrastructure/storage"
)

func TestNewS3VideoStorage_CreatesClient(t *testing.T) {
	cfg := stor.S3Config{Region: "us-east-1", Bucket: "mybucket"}
	s, err := stor.NewS3VideoStorage(cfg)
	if err != nil || s == nil {
		t.Fatalf("unexpected error creating storage: %v", err)
	}
}
