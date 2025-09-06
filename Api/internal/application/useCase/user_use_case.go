package useCase

import (
	"api/internal/domain"
	"context"
	"errors"

	"api/internal/domain/entities"
	"api/internal/domain/interfaces"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo    interfaces.UserRepository
	locRepo interfaces.LocationRepository
}

func NewUserService(repo interfaces.UserRepository, locRepo interfaces.LocationRepository) *UserService {
	return &UserService{repo: repo, locRepo: locRepo}
}

func (s *UserService) CreateUser(ctx context.Context, firstName, lastName, email, password string, country, city string) (*entities.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	cityID, err := s.locRepo.GetCityID(ctx, country, city)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: string(hash),
		CityID:       cityID,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UserService) EmailExists(ctx context.Context, email string) (bool, error) {
	_, err := s.repo.GetByEmail(ctx, email)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, domain.ErrNotFound):
		return false, nil
	default:
		return false, err
	}
}
