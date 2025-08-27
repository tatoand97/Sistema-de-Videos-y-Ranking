package ports

import "main_prj/internal/domain"

type AuthService interface {
	Register(username, email, password string, profileImagePath *string) (*domain.User, error)
	Login(email, password string) (string, *domain.User, error) // token, user
}
