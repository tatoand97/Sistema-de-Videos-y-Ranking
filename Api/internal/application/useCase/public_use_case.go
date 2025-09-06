package useCase

import (
	"context"
	"main_videork/internal/domain/interfaces"
	"main_videork/internal/domain/responses"
)

// PublicService expone operaciones públicas relacionadas con videos.
type PublicService struct {
	repo interfaces.PublicRepository
}

func NewPublicService(repo interfaces.PublicRepository) *PublicService {
	return &PublicService{repo: repo}
}

func (s *PublicService) ListPublicVideos(ctx context.Context) ([]responses.PublicVideoResponse, error) {
	return s.repo.ListPublicVideos(ctx)
}
