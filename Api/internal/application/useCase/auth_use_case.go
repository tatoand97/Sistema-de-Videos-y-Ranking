package useCase

import (
	"api/internal/domain/interfaces"
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	Permissions []string `json:"perms"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Email       string   `json:"email"`
}

type AuthService struct {
	repo          interfaces.UserRepository
	jwtSecret     string
	invalidTokens sync.Map // map[string]struct{}
}

func NewAuthService(repo interfaces.UserRepository, secret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: secret}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, int64, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", 0, errors.New("invalid credentials")
	}

	perms, err := s.repo.GetPermissions(ctx, uint(user.UserID))
	if err != nil {
		return "", 0, err
	}

	const tokenDuration = time.Hour
	now := time.Now()
	claims := AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(user.UserID),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		Permissions: perms,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", 0, err
	}
	return signed, int64(tokenDuration.Seconds()), nil
}

func (s *AuthService) Logout(token string) error {
	if token == "" {
		return errors.New("empty token")
	}
	s.invalidTokens.Store(token, struct{}{})
	return nil
}

func (s *AuthService) IsTokenInvalid(token string) bool {
	if token == "" {
		return false
	}
	_, exists := s.invalidTokens.Load(token)
	return exists
}
