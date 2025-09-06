package repository

import (
	"api/internal/domain"
	"api/internal/domain/interfaces"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
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
		// Traducir unique violation de Postgres a error de dominio (conflicto)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrConflict
		}
		return err
	}
	return nil
}
