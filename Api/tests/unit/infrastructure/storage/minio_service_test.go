package storage_test

import (
	"api/internal/domain/requests"
	stor "api/internal/infrastructure/storage"
	"context"
	"net/url"
	"testing"
)

func TestNewMinioVideoStorage_CreatesClient(t *testing.T) {
	cfg := stor.MinioConfig{Endpoint: "play.min.io", AccessKey: "minioadmin", SecretKey: "minioadmin", UseSSL: true, Bucket: "mybucket"}
	s, err := stor.NewMinioVideoStorage(cfg)
	if err != nil || s == nil {
		t.Fatalf("unexpected error creating storage: %v", err)
	}
}

func TestPresignedPostPolicy_BuildsURLAndForm(t *testing.T) {
	cfg := stor.MinioConfig{Endpoint: "play.min.io", AccessKey: "minioadmin", SecretKey: "minioadmin", UseSSL: true, Bucket: "mybucket"}
	s, err := stor.NewMinioVideoStorage(cfg)
	if err != nil || s == nil {
		t.Fatalf("unexpected error creating storage: %v", err)
	}
	// Call via interface method; should not require network
	resp, err := s.PresignedPostPolicy(context.Background(), requests.CreateUploadRequest{Filename: "test.mp4", MimeType: "video/mp4", SizeBytes: 1024})
	if err != nil {
		t.Fatalf("presigned policy error: %v", err)
	}
	if resp.UploadURL == "" || resp.Form.Policy == "" {
		t.Fatalf("expected non-empty url and policy")
	}
	if _, err := url.Parse(resp.UploadURL); err != nil {
		t.Fatalf("invalid upload URL: %v", err)
	}
}
