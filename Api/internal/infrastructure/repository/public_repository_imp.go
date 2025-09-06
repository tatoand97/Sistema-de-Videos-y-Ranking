package repository

import (
	"context"
	"main_videork/internal/domain"
	"main_videork/internal/domain/interfaces"
	"main_videork/internal/domain/responses"

	"gorm.io/gorm"
)

type publicRepository struct {
	db *gorm.DB
}

func NewPublicRepository(db *gorm.DB) interfaces.PublicRepository {
	return &publicRepository{db: db}
}

func (r *publicRepository) ListPublicVideos(ctx context.Context) ([]responses.PublicVideoResponse, error) {
	var results []responses.PublicVideoResponse
	q := r.db.WithContext(ctx).
		Table("video v").
		Select("v.video_id AS video_id, v.title, v.processed_file AS processed_url, c.name AS city, COUNT(vt.vote_id) AS votes").
		Joins("JOIN users u ON u.user_id = v.user_id").
		Joins("JOIN city c ON c.city_id = u.city_id").
		Joins("LEFT JOIN vote vt ON vt.video_id = v.video_id").
		Where("v.status = ? AND v.processed_file IS NOT NULL", "PROCESSED").
		Group("v.video_id, v.title, v.processed_file, c.name")

	if err := q.Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *publicRepository) GetPublicByID(ctx context.Context, id uint) (*responses.PublicVideoResponse, error) {
	var result responses.PublicVideoResponse
	q := r.db.WithContext(ctx).
		Table("video v").
		Select("v.video_id AS video_id, v.title, v.processed_file AS processed_url, c.name AS city, COUNT(vt.vote_id) AS votes").
		Joins("JOIN users u ON u.user_id = v.user_id").
		Joins("JOIN city c ON c.city_id = u.city_id").
		Joins("LEFT JOIN vote vt ON vt.video_id = v.video_id").
		Where("v.video_id = ? AND v.status = ? AND v.processed_file IS NOT NULL", id, "PROCESSED").
		Group("v.video_id, v.title, v.processed_file, c.name")

	if err := q.Take(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &result, nil
}
