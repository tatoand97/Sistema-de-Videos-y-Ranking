package ports

import "main_prj/internal/domain"

type UserRepository interface {
	//Create(username, passwordHash string, profileImagePath *string) (*domain.User, error)
	// TODO: Cange by email for creation
	Create(username, email, passwordHash string, profileImagePath *string) (*domain.User, error)
	//Create(email, passwordHash string) (*domain.User, error) //o quitando la imagen
	GetByUsername(username string) (*domain.User, error)
	GetByID(id int64) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
}
