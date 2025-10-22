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
	return r.CreateWithEvent(ctx, videoID, userID, nil)
}

func (r *voteRepository) CreateWithEvent(ctx context.Context, videoID, userID uint, eventID *string) error {
	type voteRow struct {
		UserID  uint    `gorm:"column:user_id"`
		VideoID uint    `gorm:"column:video_id"`
		EventID *string `gorm:"column:event_id"`
	}

	// Allow the database identity column to generate vote_id automatically.
	v := voteRow{UserID: userID, VideoID: videoID, EventID: eventID}
	if err := r.db.WithContext(ctx).Table("vote").Create(&v).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "unique_vote_user_video":
				return domain.ErrConflict
			case "ux_vote_event":
				return domain.ErrIdempotent
			default:
				// Fallback conservative: treat as conflict
				return domain.ErrConflict
			}
		}
		return err
	}
	return nil
}
