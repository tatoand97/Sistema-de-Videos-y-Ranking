package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"api/internal/domain"
	"api/internal/domain/interfaces"
)

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) interfaces.LocationRepository {
	return &locationRepository{db: db}
}

// GetCityID busca el city_id por nombre de país y ciudad (case-insensitive).
func (r *locationRepository) GetCityID(ctx context.Context, countryName, cityName string) (int, error) {
	if strings.TrimSpace(countryName) == "" || strings.TrimSpace(cityName) == "" {
		return 0, domain.ErrInvalid
	}

	var id int
	q := r.db.WithContext(ctx).
		Table("city c").
		Select("c.city_id").
		Joins("JOIN country co ON co.country_id = c.country_id").
		Where("LOWER(co.name) = ? AND LOWER(c.name) = ?", strings.ToLower(countryName), strings.ToLower(cityName)).
		Limit(1)

	if err := q.Scan(&id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, domain.ErrNotFound
		}
		return 0, err
	}
	if id == 0 {
		return 0, domain.ErrNotFound
	}
	return id, nil
}
