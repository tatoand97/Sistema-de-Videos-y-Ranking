package useCase

import (
	"api/internal/domain"
	"api/internal/domain/interfaces"
	"api/internal/domain/responses"
	"context"
	"errors"
	"strconv"
	"strings"
)

// PublicService expone operaciones publicas relacionadas con videos.
type PublicService struct {
	repo     interfaces.PublicRepository
	voteRepo interfaces.VoteRepository
	agg      interfaces.Aggregates
}

func NewPublicService(repo interfaces.PublicRepository, voteRepo interfaces.VoteRepository) *PublicService {
	return &PublicService{repo: repo, voteRepo: voteRepo}
}

// NewPublicServiceWithAgg permite inyectar un lector de agregados (Redis) sin acoplar al handler.
func NewPublicServiceWithAgg(repo interfaces.PublicRepository, voteRepo interfaces.VoteRepository, agg interfaces.Aggregates) *PublicService {
	return &PublicService{repo: repo, voteRepo: voteRepo, agg: agg}
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
	if s.agg != nil {
		// Construir pollId de usuarios (global o por ciudad)
		poll := "users"
		if city != nil && *city != "" {
			poll = poll + ":city:" + normalizeCityKey(*city)
		}
		start := int64((page - 1) * pageSize)
		stop := start + int64(pageSize) - 1
		items, err := s.agg.GetLeaderboard(ctx, poll, start, stop)
		if err == nil && len(items) > 0 {
			// Enriquecer con username y ciudad desde BD
			ids := make([]uint, 0, len(items))
			for _, it := range items {
				if uid64, e := strconv.ParseUint(it.Member, 10, 64); e == nil {
					ids = append(ids, uint(uid64))
				}
			}
			basics, berr := s.repo.GetUsersBasicByIDs(ctx, ids)
			if berr == nil {
				bm := make(map[uint]responses.UserBasic, len(basics))
				for _, b := range basics {
					bm[b.UserID] = b
				}
				out := make([]responses.RankingItem, 0, len(items))
				for _, it := range items {
					uid64, _ := strconv.ParseUint(it.Member, 10, 64)
					b := bm[uint(uid64)]
					out = append(out, responses.RankingItem{Username: b.Username, City: b.City, Votes: int(it.Score)})
				}
				return out, nil
			}
		}
		// Si no hay datos en agregados o hay error -> caer a BD
	}
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
func (s *PublicService) GetUsersBasicByIDs(ctx context.Context, ids []uint) ([]responses.UserBasic, error) {
	return s.repo.GetUsersBasicByIDs(ctx, ids)
}
