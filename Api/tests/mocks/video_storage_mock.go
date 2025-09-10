package mocks

import (
	"context"
	"io"
)

type MockVideoStorage struct {
	SaveFunc func(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
}

func (m *MockVideoStorage) Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, objectName, reader, size, contentType)
	}
	return "https://example.com/video.mp4", nil
}
