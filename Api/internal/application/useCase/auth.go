package useCase

import (
	"context"
	"errors"
	"main_videork/internal/domain"
	"main_videork/internal/domain/entities"
	"main_videork/internal/domain/interfaces"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo          interfaces.UserRepository
	jwtSecret     string
	invalidTokens sync.Map // map[string]struct{}
}

func NewAuthService(repo interfaces.UserRepository, secret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: secret}
}

func (s *AuthService) Register(ctx context.Context, firstName, lastName, email, password, city, country string) (*entities.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: string(hash),
		City:         city,
		Country:      country,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, int64, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", 0, errors.New("invalid credentials")
	}

	const tokenDuration = time.Hour
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(user.UserId), 10),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenDuration)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", 0, err
	}
	return signed, int64(tokenDuration.Seconds()), nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return errors.New("empty token")
	}
	s.invalidTokens.Store(token, struct{}{})
	return nil
}

func (s *AuthService) EmailExists(ctx context.Context, email string) (bool, error) {
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
