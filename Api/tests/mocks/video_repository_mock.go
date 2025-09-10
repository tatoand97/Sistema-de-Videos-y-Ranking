package mocks

import (
	"api/internal/domain/entities"
	"context"
)

type MockVideoRepository struct {
	CreateFunc         func(ctx context.Context, video *entities.Video) error
	GetByIDFunc        func(ctx context.Context, id uint) (*entities.Video, error)
	ListFunc           func(ctx context.Context) ([]*entities.Video, error)
	ListByUserFunc     func(ctx context.Context, userID uint) ([]*entities.Video, error)
	GetByIDAndUserFunc func(ctx context.Context, id, userID uint) (*entities.Video, error)
	DeleteFunc         func(ctx context.Context, id uint) error
	UpdateStatusFunc   func(ctx context.Context, id uint, status entities.VideoStatus) error
}

func (m *MockVideoRepository) Create(ctx context.Context, video *entities.Video) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, video)
	}
	return nil
}

func (m *MockVideoRepository) GetByID(ctx context.Context, id uint) (*entities.Video, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockVideoRepository) List(ctx context.Context) ([]*entities.Video, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}

func (m *MockVideoRepository) ListByUser(ctx context.Context, userID uint) ([]*entities.Video, error) {
	if m.ListByUserFunc != nil {
		return m.ListByUserFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockVideoRepository) GetByIDAndUser(ctx context.Context, id, userID uint) (*entities.Video, error) {
	if m.GetByIDAndUserFunc != nil {
		return m.GetByIDAndUserFunc(ctx, id, userID)
	}
	return nil, nil
}

func (m *MockVideoRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockVideoRepository) UpdateStatus(ctx context.Context, id uint, status entities.VideoStatus) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil
}
