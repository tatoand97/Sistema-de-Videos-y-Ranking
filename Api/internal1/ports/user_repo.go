package ports

import "main_viderk/internal/domain"

type UserRepository interface {
	Create(username, email, passwordHash string, profileImagePath *string) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	GetByID(id int64) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
}
