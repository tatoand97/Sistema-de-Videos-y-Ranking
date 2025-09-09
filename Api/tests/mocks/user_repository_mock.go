package mocks

import (
	"api/internal/domain/entities"
	"context"
)

type MockUserRepository struct {
	GetByEmailFunc     func(ctx context.Context, email string) (*entities.User, error)
	GetPermissionsFunc func(ctx context.Context, userID uint) ([]string, error)
	CreateFunc         func(ctx context.Context, user *entities.User) error
	GetByIDFunc        func(ctx context.Context, id uint) (*entities.User, error)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) GetPermissions(ctx context.Context, userID uint) ([]string, error) {
	if m.GetPermissionsFunc != nil {
		return m.GetPermissionsFunc(ctx, userID)
	}
	return []string{}, nil
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*entities.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}
