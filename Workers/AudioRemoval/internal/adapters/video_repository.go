package adapters

import (
	"audioremoval/internal/domain"
	"time"
)

type VideoRepository struct{}

func NewVideoRepository() *VideoRepository {
	return &VideoRepository{}
}

func (r *VideoRepository) FindByFilename(filename string) (*domain.Video, error) {
	return &domain.Video{
		ID:        filename,
		Filename:  filename,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}, nil
}

func (r *VideoRepository) UpdateStatus(id string, status domain.ProcessingStatus) error {
	// En producción: actualizar en base de datos
	return nil
}