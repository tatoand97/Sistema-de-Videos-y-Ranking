package adapters

import (
	"statesmachine/internal/domain"
	"fmt"
)

type MockVideoRepository struct{}

func NewMockVideoRepository() *MockVideoRepository {
	return &MockVideoRepository{}
}

func (r *MockVideoRepository) FindByFilename(filename string) (*domain.Video, error) {
	return &domain.Video{
		ID:       fmt.Sprintf("video_%s", filename),
		Filename: filename,
		Status:   domain.StatusPending,
	}, nil
}

func (r *MockVideoRepository) UpdateStatus(id string, status domain.VideoStatus) error {
	return nil
}