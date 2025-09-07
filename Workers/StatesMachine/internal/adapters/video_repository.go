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
	updates := map[string]interface{}{
		"status": string(status),
	}
	
	// Handle processed_at according to constraint: only set when status is PROCESSED
	if status == domain.StatusProcessed {
		updates["processed_at"] = "NOW()"
	} else {
		updates["processed_at"] = nil
	}
	
	return r.db.Model(&domain.Video{}).Where("video_id = ?", id).Updates(updates).Error
}