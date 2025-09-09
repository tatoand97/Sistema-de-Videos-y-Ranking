package application_test

import (
	"api/internal/application/useCase"
	"api/internal/domain/entities"
	"api/tests/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name         string
		email        string
		password     string
		mockRepo     *mocks.MockUserRepository
		wantErr      bool
		wantToken    bool
		wantDuration bool
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "password123",
			mockRepo: &mocks.MockUserRepository{
				GetByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
					return &entities.User{
						UserID:       1,
						Email:        email,
						PasswordHash: string(hashedPassword),
						FirstName:    "John",
						LastName:     "Doe",
					}, nil
				},
				GetPermissionsFunc: func(ctx context.Context, userID uint) ([]string, error) {
					return []string{"read", "write"}, nil
				},
			},
			wantErr:      false,
			wantToken:    true,
			wantDuration: true,
		},
		{
			name:     "user not found",
			email:    "notfound@example.com",
			password: "password123",
			mockRepo: &mocks.MockUserRepository{
				GetByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
					return nil, errors.New("user not found")
				},
			},
			wantErr: true,
		},
		{
			name:     "invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			mockRepo: &mocks.MockUserRepository{
				GetByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
					return &entities.User{
						UserID:       1,
						Email:        email,
						PasswordHash: string(hashedPassword),
						FirstName:    "John",
						LastName:     "Doe",
					}, nil
				},
			},
			wantErr: true,
		},
		{
			name:     "permissions error",
			email:    "test@example.com",
			password: "password123",
			mockRepo: &mocks.MockUserRepository{
				GetByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
					return &entities.User{
						UserID:       1,
						Email:        email,
						PasswordHash: string(hashedPassword),
						FirstName:    "John",
						LastName:     "Doe",
					}, nil
				},
				GetPermissionsFunc: func(ctx context.Context, userID uint) ([]string, error) {
					return nil, errors.New("permissions error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := useCase.NewAuthService(tt.mockRepo, "test-secret")

			token, duration, err := authService.Login(context.Background(), tt.email, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.Zero(t, duration)
			} else {
				assert.NoError(t, err)
				if tt.wantToken {
					assert.NotEmpty(t, token)
				}
				if tt.wantDuration {
					assert.Greater(t, duration, int64(0))
				}
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "successful logout",
			token:   "valid-token",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := useCase.NewAuthService(nil, "test-secret")

			err := authService.Logout(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify token is marked as invalid
				assert.True(t, authService.IsTokenInvalid(tt.token))
			}
		})
	}
}

func TestAuthService_IsTokenInvalid(t *testing.T) {
	authService := useCase.NewAuthService(nil, "test-secret")

	// Initially, token should not be invalid
	assert.False(t, authService.IsTokenInvalid("test-token"))

	// After logout, token should be invalid
	authService.Logout("test-token")
	assert.True(t, authService.IsTokenInvalid("test-token"))

	// Empty token should return false
	assert.False(t, authService.IsTokenInvalid(""))
}
