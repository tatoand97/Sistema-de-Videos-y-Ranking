package interfaces

import (
	"context"
	"io"
)

// VideoStorage defines behavior for storing video files.
type VideoStorage interface {
	// Save stores a video object with the given name and returns the object key or URL.
	Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
}
