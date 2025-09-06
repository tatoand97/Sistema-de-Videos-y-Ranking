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

func NewPublicService(repo interfaces.PublicRepository) *PublicService { // backward compat
	return &PublicService{repo: repo}
}

func (s *PublicService) ListPublicVideos(ctx context.Context) ([]responses.PublicVideoResponse, error) {
	return s.repo.ListPublicVideos(ctx)
}

// WithVotes permite inyectar el repositorio de votos sin romper firmas existentes
func (s *PublicService) WithVotes(voteRepo interfaces.VoteRepository) *PublicService {
	s.voteRepo = voteRepo
	return s
}

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
