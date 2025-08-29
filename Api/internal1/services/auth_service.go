package services

import (
	"errors"
	"time"

	"main_viderk/internal/domain"
	"main_viderk/internal/ports"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	users     ports.UserRepository
	jwtSecret string
}

func NewAuthService(users ports.UserRepository, jwtSecret string) ports.AuthService {
	return &authService{users: users, jwtSecret: jwtSecret}
}

func (s *authService) Register(username, email, password string, profileImagePath *string) (*domain.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email/password required")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return s.users.Create(username, email, string(hash), profileImagePath)
}

func (s *authService) Login(email, password string) (string, *domain.User, error) {
	u, err := s.users.GetByEmail(email)
	if err != nil || u == nil {
		return "", nil, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", nil, errors.New("invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      u.ID,
		"username": u.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	sToken, err := token.SignedString([]byte(s.jwtSecret))
	return sToken, u, err
}
