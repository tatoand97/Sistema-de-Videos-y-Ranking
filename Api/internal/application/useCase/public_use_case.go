package useCase

import (
	"api/internal/domain"
	"api/internal/domain/interfaces"
	"api/internal/domain/responses"
	"context"
	"errors"
)

// PublicService expone operaciones públicas relacionadas con videos.
type PublicService struct {
	repo     interfaces.PublicRepository
	voteRepo interfaces.VoteRepository
}

func NewPublicService(repo interfaces.PublicRepository, voteRepo interfaces.VoteRepository) *PublicService {
	return &PublicService{repo: repo, voteRepo: voteRepo}
}

func (s *PublicService) ListPublicVideos(ctx context.Context) ([]responses.PublicVideoResponse, error) {
	return s.repo.ListPublicVideos(ctx)
}

// WithVotes was removed; voteRepo is injected in constructor.

func (s *PublicService) GetPublicByID(ctx context.Context, id uint) (*responses.PublicVideoResponse, error) {
	return s.repo.GetPublicByID(ctx, id)
}

// VotePublicVideo aplica la regla de un voto por usuario por video
func (s *PublicService) VotePublicVideo(ctx context.Context, videoID, userID uint) error {
	if s.voteRepo == nil {
		return errors.New("vote repository not configured")
	}
	// Verificamos existencia y que sea público
	if _, err := s.repo.GetPublicByID(ctx, videoID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	// Verifica si ya votó
	already, err := s.voteRepo.HasUserVoted(ctx, videoID, userID)
	if err != nil {
		return err
	}
	if already {
		return domain.ErrConflict
	}
	// Crear voto (con índice único anti-race)
	if err := s.voteRepo.Create(ctx, videoID, userID); err != nil {
		return err
	}
	return nil
}

// Rankings retorna el ranking paginado por votos acumulados por usuario.
func (s *PublicService) Rankings(ctx context.Context, city *string, page, pageSize int) ([]responses.RankingItem, error) {
	return s.repo.Rankings(ctx, city, page, pageSize)
}
