package useCase

import (
	"api/internal/domain"
	"api/internal/domain/interfaces"
	"api/internal/domain/responses"
	"context"
	"errors"
	"strings"
)

// PublicService expone operaciones publicas relacionadas con videos.
type PublicService struct {
	repo     interfaces.PublicRepository
	voteRepo interfaces.VoteRepository
}

func NewPublicService(repo interfaces.PublicRepository, voteRepo interfaces.VoteRepository) *PublicService {
	return &PublicService{repo: repo, voteRepo: voteRepo}
}

// NewPublicServiceWithAgg permite inyectar un lector de agregados (Redis) sin acoplar al handler.
// NewPublicServiceWithAgg eliminado (sin agregados Redis)

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
	// Verificamos existencia y que sea publico
	if _, err := s.repo.GetPublicByID(ctx, videoID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	// Verifica si ya voto
	already, err := s.voteRepo.HasUserVoted(ctx, videoID, userID)
	if err != nil {
		return err
	}
	if already {
		return domain.ErrConflict
	}
	// Crear voto (con indice unico anti-race)
	if err := s.voteRepo.Create(ctx, videoID, userID); err != nil {
		return err
	}
	return nil
}

// VotePublicVideoWithEvent inserta el voto persistiendo eventID para idempotencia fuerte.
// En caso de duplicado por eventID devuelve domain.ErrIdempotent.
func (s *PublicService) VotePublicVideoWithEvent(ctx context.Context, videoID, userID uint, eventID *string) error {
	if s.voteRepo == nil {
		return errors.New("vote repository not configured")
	}
	// Verificar existencia de video publico
	if _, err := s.repo.GetPublicByID(ctx, videoID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	// Insertar con eventID si el repo lo soporta, si no fallback a Create
	if repoEvt, ok := interface{}(s.voteRepo).(interfaces.VoteRepositoryWithEvent); ok {
		if err := repoEvt.CreateWithEvent(ctx, videoID, userID, eventID); err != nil {
			if errors.Is(err, domain.ErrIdempotent) {
				return domain.ErrIdempotent
			}
			return err
		}
	} else {
		if err := s.voteRepo.Create(ctx, videoID, userID); err != nil {
			return err
		}
	}
	return nil
}

// Rankings retorna el ranking paginado por votos acumulados por usuario.
func (s *PublicService) Rankings(ctx context.Context, city *string, page, pageSize int) ([]responses.RankingItem, error) {
	// Orquestación: primero intentar desde agregados (Redis) si está disponible.
	return s.repo.Rankings(ctx, city, page, pageSize)
}

// normalizeCityKey replica la normalización usada en handlers para claves por ciudad.
func normalizeCityKey(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return s
	}
	r := strings.NewReplacer(
		"á", "a", "à", "a", "ä", "a", "â", "a", "ã", "a",
		"é", "e", "è", "e", "ë", "e", "ê", "e",
		"í", "i", "ì", "i", "ï", "i", "î", "i",
		"ó", "o", "ò", "o", "ö", "o", "ô", "o", "õ", "o",
		"ú", "u", "ù", "u", "ü", "u", "û", "u",
		"ñ", "n",
	)
	s = r.Replace(s)
	s = strings.ReplaceAll(s, " ", "-")
	b := make([]rune, 0, len(s))
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' {
			b = append(b, ch)
		}
	}
	return string(b)
}

// GetUsersBasicByIDs expone informacion basica de usuarios para enriquecer rankings desde Redis.
// GetUsersBasicByIDs removido; se usaba solo para enriquecer rankings desde Redis.
