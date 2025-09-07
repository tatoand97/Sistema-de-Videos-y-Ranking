package adapters

import (
	"statesmachine/internal/domain"
	"errors"
	"gorm.io/gorm"
)

type PostgresVideoRepository struct {
	db *gorm.DB
}

func NewPostgresVideoRepository(db *gorm.DB) *PostgresVideoRepository {
	return &PostgresVideoRepository{db: db}
}

func (r *PostgresVideoRepository) FindByID(id uint) (*domain.Video, error) {
	var video domain.Video
	if err := r.db.First(&video, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("video not found")
		}
		return nil, err
	}
	return &video, nil
}

func (r *PostgresVideoRepository) UpdateStatus(id uint, status domain.VideoStatus) error {
	return r.db.Model(&domain.Video{}).Where("video_id = ?", id).Update("status", string(status)).Error
}