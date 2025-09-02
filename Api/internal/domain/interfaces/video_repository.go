package interfaces

import (
	"context"
	"main_videork/internal/domain/entities"
)

// VideoRepository define el comportamiento de persistencia para videos.
type VideoRepository interface {
	Create(ctx context.Context, video *entities.Video) error
	GetByID(ctx context.Context, id uint) (*entities.Video, error)
	List(ctx context.Context) ([]*entities.Video, error)
}
