package repository

import (
	"api/internal/domain"
	"api/internal/domain/interfaces"
	"api/internal/domain/responses"
	"context"

	"gorm.io/gorm"
)

const joinCityOnUser = "JOIN city c ON c.city_id = u.city_id"
const leftJoinVoteOnVideo = "LEFT JOIN vote vt ON vt.video_id = v.video_id"

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
		Select("v.video_id AS video_id, v.title, CONCAT('http://localhost:8081/processed-videos/', v.processed_file) AS processed_url, c.name AS city, COUNT(vt.vote_id) AS votes").
		Joins("JOIN users u ON u.user_id = v.user_id").
		Joins(joinCityOnUser).
		Joins(leftJoinVoteOnVideo).
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
		Select("v.video_id AS video_id, v.title, CONCAT('http://localhost:8081/processed-videos/', v.processed_file) AS processed_url, c.name AS city, COUNT(vt.vote_id) AS votes").
		Joins("JOIN users u ON u.user_id = v.user_id").
		Joins(joinCityOnUser).
		Joins(leftJoinVoteOnVideo).
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

// Rankings agrega votos por usuario (solo sobre videos públicos procesados) y ordena por votos desc.
// Desempate TBD: se aplica orden estable por user_id (no expuesto).
func (r *publicRepository) Rankings(ctx context.Context, city *string, page, pageSize int) ([]responses.RankingItem, error) {
	type row struct {
		Username string  `gorm:"column:username"`
		City     *string `gorm:"column:city"`
		Votes    int     `gorm:"column:votes"`
	}

	offset := (page - 1) * pageSize

	q := r.db.WithContext(ctx).
		Table("users u").
		Select("split_part(u.email, '@', 1) AS username, c.name AS city, COUNT(vt.vote_id) AS votes").
		Joins(joinCityOnUser).
		Joins("JOIN video v ON v.user_id = u.user_id").
		Joins(leftJoinVoteOnVideo).
		Where("v.status = ? AND v.processed_file IS NOT NULL", "PROCESSED").
		Group("u.user_id, u.email, c.name").
		Order("votes DESC, u.user_id ASC"). // desempate interno estable (TBD)
		Limit(pageSize).
		Offset(offset)

	if city != nil && *city != "" {
		// Filtro de ciudad sin tildes ni mayúsculas/minúsculas, usando wrapper inmutable para que el planificador use el índice
		q = q.Where("immutable_unaccent(LOWER(c.name)) = immutable_unaccent(LOWER(?))", *city)
	}

	var rows []row
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]responses.RankingItem, 0, len(rows))
	for _, rrow := range rows {
		items = append(items, responses.RankingItem{
			Username: rrow.Username,
			City:     rrow.City,
			Votes:    rrow.Votes,
		})
	}
	return items, nil
}
