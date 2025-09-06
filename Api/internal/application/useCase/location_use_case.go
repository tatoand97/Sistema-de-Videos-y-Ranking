package useCase

import (
	"api/internal/domain/interfaces"
	"context"
)

// LocationService expone operaciones de lectura para ubicación (país/ciudad).
type LocationService struct {
	repo interfaces.LocationRepository
}

func NewLocationService(repo interfaces.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

// GetCityID retorna el city_id dado país y ciudad (case-insensitive según repo).
func (s *LocationService) GetCityID(ctx context.Context, country, city string) (int, error) {
	return s.repo.GetCityID(ctx, country, city)
}
