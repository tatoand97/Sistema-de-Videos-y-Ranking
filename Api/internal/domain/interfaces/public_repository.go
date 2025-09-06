package interfaces

import (
	"api/internal/domain/responses"
	"context"
)

// PublicRepository define operaciones de lectura públicas (agregadas) para videos.
type PublicRepository interface {
	ListPublicVideos(ctx context.Context) ([]responses.PublicVideoResponse, error)
	// GetPublicByID devuelve el video público (procesado) por ID o ErrNotFound
	GetPublicByID(ctx context.Context, id uint) (*responses.PublicVideoResponse, error)
	// Rankings devuelve un listado paginado del ranking de jugadores por votos acumulados.
	// city: filtro opcional por nombre de ciudad (case-insensitive). Si nil, no filtra.
	// page, pageSize: para paginación.
	Rankings(ctx context.Context, city *string, page, pageSize int) ([]responses.RankingItem, error)
}
