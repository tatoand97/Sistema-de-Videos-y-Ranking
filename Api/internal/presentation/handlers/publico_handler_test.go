package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"main_videork/internal/application/useCase"
	"main_videork/internal/domain"
	"main_videork/internal/domain/responses"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

// Fake implementations of repositories
type fakePublicRepo struct{ exists bool }

func (f *fakePublicRepo) ListPublicVideos(ctx context.Context) ([]responses.PublicVideoResponse, error) {
	return []responses.PublicVideoResponse{}, nil
}
func (f *fakePublicRepo) GetPublicByID(ctx context.Context, id uint) (*responses.PublicVideoResponse, error) {
	if f.exists {
		r := responses.PublicVideoResponse{VideoID: id, Title: "ok"}
		return &r, nil
	}
	return nil, domain.ErrNotFound
}

type fakeVoteRepo struct {
	hasVoted  bool
	createErr error
}

func (f *fakeVoteRepo) HasUserVoted(ctx context.Context, videoID, userID uint) (bool, error) {
	return f.hasVoted, nil
}
func (f *fakeVoteRepo) Create(ctx context.Context, videoID, userID uint) error { return f.createErr }

// Helper to setup Gin with handler and a middleware to set userID
func setupRouter(h *PublicHandlers, withAuth bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if withAuth {
		r.Use(func(c *gin.Context) { c.Set("userID", uint(10)); c.Next() })
	}
	r.POST("/api/public/videos/:video_id/vote", h.VotePublicVideo)
	return r
}

func TestVotePublicVideo_OK(t *testing.T) {
	pub := &fakePublicRepo{exists: true}
	votes := &fakeVoteRepo{hasVoted: false, createErr: nil}
	svc := useCase.NewPublicService(pub).WithVotes(votes)
	h := NewPublicHandlers(svc)
	r := setupRouter(h, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/public/videos/123/vote", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVotePublicVideo_AlreadyVoted(t *testing.T) {
	pub := &fakePublicRepo{exists: true}
	votes := &fakeVoteRepo{hasVoted: true}
	svc := useCase.NewPublicService(pub).WithVotes(votes)
	h := NewPublicHandlers(svc)
	r := setupRouter(h, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/public/videos/123/vote", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVotePublicVideo_UniqueViolation(t *testing.T) {
	pub := &fakePublicRepo{exists: true}
	// simulate unique violation 23505
	pgErr := &pgconn.PgError{Code: "23505"}
	votes := &fakeVoteRepo{hasVoted: false, createErr: pgErr}
	svc := useCase.NewPublicService(pub).WithVotes(votes)
	h := NewPublicHandlers(svc)
	r := setupRouter(h, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/public/videos/123/vote", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVotePublicVideo_NotFound(t *testing.T) {
	pub := &fakePublicRepo{exists: false}
	votes := &fakeVoteRepo{}
	svc := useCase.NewPublicService(pub).WithVotes(votes)
	h := NewPublicHandlers(svc)
	r := setupRouter(h, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/public/videos/123/vote", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVotePublicVideo_Unauthorized(t *testing.T) {
	pub := &fakePublicRepo{exists: true}
	votes := &fakeVoteRepo{}
	svc := useCase.NewPublicService(pub).WithVotes(votes)
	h := NewPublicHandlers(svc)
	r := setupRouter(h, false) // no auth middleware

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/public/videos/123/vote", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body=%s", w.Code, w.Body.String())
	}
}
