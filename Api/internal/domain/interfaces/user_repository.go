package interfaces

import (
	"context"
	"main_videork/internal/domain/entities"
)

// UserRepository defines persistence behavior for users.
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
}
