package interfaces

import (
	"context"
	"main_videork/internal/domain/entities"
)

// UserRepository defines persistence behavior for users.
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	// GetPermissions returns all privilege names assigned to the user through roles.
	GetPermissions(ctx context.Context, userID uint) ([]string, error)
}
