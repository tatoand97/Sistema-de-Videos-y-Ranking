package interfaces

import (
	"context"
	"main_videork/internal/domain/responses"
)

// PublicRepository define operaciones de lectura públicas (agregadas) para videos.
type PublicRepository interface {
	ListPublicVideos(ctx context.Context) ([]responses.PublicVideoResponse, error)
	// GetPublicByID devuelve el video público (procesado) por ID o ErrNotFound
	GetPublicByID(ctx context.Context, id uint) (*responses.PublicVideoResponse, error)
}
