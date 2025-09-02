package adapters

import (
	"gossipopenclose/internal/domain"
	"time"
)

type VideoRepository struct{}

func NewVideoRepository() *VideoRepository { return &VideoRepository{} }

func (r *VideoRepository) FindByFilename(filename string) (*domain.Video, error) {
	return &domain.Video{
		ID: filename,
		Filename: filename,
		Status: domain.StatusPending,
		CreatedAt: time.Now(),
	}, nil
}

func (r *VideoRepository) UpdateStatus(id string, status domain.ProcessingStatus) error {
	// En producción: persistir estado (DB)
	return nil
}
