package interfaces

import (
	"api/internal/domain/responses"
	"context"
)

// PublicRepository define operaciones de lectura publicas (agregadas) para videos.
type PublicRepository interface {
	ListPublicVideos(ctx context.Context) ([]responses.PublicVideoResponse, error)
	// GetPublicByID devuelve el video publico (procesado) por ID o ErrNotFound
	GetPublicByID(ctx context.Context, id uint) (*responses.PublicVideoResponse, error)
	// Rankings devuelve un listado paginado del ranking de jugadores por votos acumulados.
	// city: filtro opcional por nombre de ciudad (case-insensitive). Si nil, no filtra.
	// page, pageSize: para paginacion.
	Rankings(ctx context.Context, city *string, page, pageSize int) ([]responses.RankingItem, error)
	// GetUsersBasicByIDs retorna username y ciudad para los userIDs proporcionados.
	GetUsersBasicByIDs(ctx context.Context, ids []uint) ([]responses.UserBasic, error)
}
