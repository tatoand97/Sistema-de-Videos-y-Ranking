package repository

import (
	"api/internal/domain"
	"api/internal/domain/entities"
	"api/internal/domain/interfaces"
	"context"
	"errors"

	"gorm.io/gorm"
)

type videoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) interfaces.VideoRepository {
	return &videoRepository{db: db}
}

func (r *videoRepository) Create(ctx context.Context, video *entities.Video) error {
	return r.db.WithContext(ctx).Create(video).Error
}

func (r *videoRepository) GetByID(ctx context.Context, id uint) (*entities.Video, error) {
	var video entities.Video
	if err := r.db.WithContext(ctx).First(&video, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) List(ctx context.Context) ([]*entities.Video, error) {
	var videos []*entities.Video
	if err := r.db.WithContext(ctx).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func (r *videoRepository) ListByUser(ctx context.Context, userID uint) ([]*entities.Video, error) {
	var videos []*entities.Video
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("uploaded_at DESC").
		Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// GetByIDAndUser fetches a video by id and ensures it belongs to the given user.
// Distinguishes between not found (404) and forbidden (403).
func (r *videoRepository) GetByIDAndUser(ctx context.Context, id, userID uint) (*entities.Video, error) {
	// First, try to get by id to know if it exists at all
	var v entities.Video
	if err := r.db.WithContext(ctx).First(&v, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	if v.UserID != userID {
		return nil, domain.ErrForbidden
	}
	return &v, nil
}

func (r *videoRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&entities.Video{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
