package mocks

import (
	"api/internal/domain/requests"
	"api/internal/domain/responses"
	"context"
	"io"
)

type MockVideoStorage struct {
	SaveFunc                func(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
	PresignedPostPolicyFunc func(ctx context.Context, req requests.CreateUploadRequest) (*responses.CreateUploadResponsePostPolicy, error)
}

func (m *MockVideoStorage) Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, objectName, reader, size, contentType)
	}
	return "https://example.com/video.mp4", nil
}

func (m *MockVideoStorage) PresignedPostPolicy(ctx context.Context, req requests.CreateUploadRequest) (*responses.CreateUploadResponsePostPolicy, error) {
	if m.PresignedPostPolicyFunc != nil {
		return m.PresignedPostPolicyFunc(ctx, req)
	}
	return &responses.CreateUploadResponsePostPolicy{}, nil
}