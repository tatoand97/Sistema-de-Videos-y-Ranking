package interfaces

import (
	"api/internal/domain/requests"
	"api/internal/domain/responses"
	"context"
	"io"
)

// VideoStorage defines behavior for storing video files.
type VideoStorage interface {
	// Save stores a video object with the given name and returns the object key or URL.
	Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
	// PresignedPostPolicy returns a signed S3 POST policy to upload directly to storage.
	PresignedPostPolicy(ctx context.Context, req requests.CreateUploadRequest) (*responses.CreateUploadResponsePostPolicy, error)
}
