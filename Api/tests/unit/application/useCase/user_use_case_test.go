package application_test

import (
	usecase "api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/entities"
	"api/tests/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockLocRepo struct {
	GetCityIDFunc func(ctx context.Context, countryName, cityName string) (int, error)
}

func (m *mockLocRepo) GetCityID(ctx context.Context, countryName, cityName string) (int, error) {
	if m.GetCityIDFunc != nil {
		return m.GetCityIDFunc(ctx, countryName, cityName)
	}
	return 0, nil
}

func TestUserService_CreateUser_Success(t *testing.T) {
	var captured *entities.User
	userRepo := &mocks.MockUserRepository{
		CreateFunc: func(ctx context.Context, u *entities.User) error { captured = u; return nil },
	}
	locRepo := &mockLocRepo{GetCityIDFunc: func(ctx context.Context, country, city string) (int, error) {
		assert.Equal(t, "Peru", country)
		assert.Equal(t, "Lima", city)
		return 77, nil
	}}
	svc := usecase.NewUserService(userRepo, locRepo)

	u, err := svc.CreateUser(context.Background(), "Ana", "Lopez", "ana@example.com", "s3cret", "Peru", "Lima")
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, 77, u.CityID)
	assert.Equal(t, "ana@example.com", u.Email)
	assert.NotEmpty(t, u.PasswordHash)
	assert.NotEqual(t, "s3cret", u.PasswordHash)
	// Ensure repo received the same user pointer
	assert.Same(t, captured, u)
}

func TestUserService_CreateUser_LocationError(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	locRepo := &mockLocRepo{GetCityIDFunc: func(ctx context.Context, country, city string) (int, error) {
		return 0, errors.New("loc fail")
	}}
	svc := usecase.NewUserService(userRepo, locRepo)

	u, err := svc.CreateUser(context.Background(), "Ana", "Lopez", "ana@example.com", "pw", "X", "Y")
	assert.Error(t, err)
	assert.Nil(t, u)
}

func TestUserService_CreateUser_RepoError(t *testing.T) {
	userRepo := &mocks.MockUserRepository{CreateFunc: func(ctx context.Context, u *entities.User) error { return errors.New("db fail") }}
	locRepo := &mockLocRepo{GetCityIDFunc: func(ctx context.Context, country, city string) (int, error) { return 50, nil }}
	svc := usecase.NewUserService(userRepo, locRepo)

	u, err := svc.CreateUser(context.Background(), "Ana", "Lopez", "ana@example.com", "pw", "Peru", "Lima")
	assert.Error(t, err)
	assert.Nil(t, u)
}

func TestUserService_GetByEmail(t *testing.T) {
	userRepo := &mocks.MockUserRepository{GetByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
		return &entities.User{Email: email}, nil
	}}
	locRepo := &mockLocRepo{}
	svc := usecase.NewUserService(userRepo, locRepo)

	u, err := svc.GetByEmail(context.Background(), "ana@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "ana@example.com", u.Email)
}

func TestUserService_EmailExists(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{GetByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
			return &entities.User{Email: email}, nil
		}}
		svc := usecase.NewUserService(userRepo, &mockLocRepo{})
		exists, err := svc.EmailExists(context.Background(), "a@b.com")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("not found", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{GetByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
			return nil, domain.ErrNotFound
		}}
		svc := usecase.NewUserService(userRepo, &mockLocRepo{})
		exists, err := svc.EmailExists(context.Background(), "a@b.com")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("repo error", func(t *testing.T) {
		userRepo := &mocks.MockUserRepository{GetByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
			return nil, errors.New("db down")
		}}
		svc := usecase.NewUserService(userRepo, &mockLocRepo{})
		exists, err := svc.EmailExists(context.Background(), "a@b.com")
		assert.Error(t, err)
		assert.False(t, exists)
	})
}
