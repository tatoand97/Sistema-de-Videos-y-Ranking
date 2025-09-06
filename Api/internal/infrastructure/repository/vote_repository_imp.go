package repository

import (
	"context"
	"main_videork/internal/domain/interfaces"

	"gorm.io/gorm"
)

type voteRepository struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) interfaces.VoteRepository {
	return &voteRepository{db: db}
}

func (r *voteRepository) HasUserVoted(ctx context.Context, videoID, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Table("vote").
		Where("video_id = ? AND user_id = ?", videoID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *voteRepository) Create(ctx context.Context, videoID, userID uint) error {
	// Insert directo con tabla "vote" para simplicidad.
	type voteRow struct {
		VoteID  uint `gorm:"column:vote_id"`
		UserID  uint `gorm:"column:user_id"`
		VideoID uint `gorm:"column:video_id"`
	}
	v := voteRow{UserID: userID, VideoID: videoID}
	if err := r.db.WithContext(ctx).Table("vote").Create(&v).Error; err != nil {
		// Propaga el error tal cual (puede ser unique violation 23505)
		return err
	}
	return nil
}
