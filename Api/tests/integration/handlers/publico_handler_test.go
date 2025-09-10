package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/responses"
	hdl "api/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

// Fake implementations of repositories
type fakePublicRepo struct {
	exists   bool
	rankings []responses.RankingItem
}

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
func (f *fakePublicRepo) Rankings(ctx context.Context, city *string, page, pageSize int) ([]responses.RankingItem, error) {
	// Filter by city if provided (case-insensitive exact)
	var filtered []responses.RankingItem
	if city != nil && *city != "" {
		cLower := strings.ToLower(*city)
		for _, it := range f.rankings {
			if it.City == nil {
				continue
			}
			if strings.ToLower(*it.City) == cLower {
				filtered = append(filtered, it)
			}
		}
	} else {
		filtered = append(filtered, f.rankings...)
	}
	// Pagination
	start := (page - 1) * pageSize
	if start >= len(filtered) {
		return []responses.RankingItem{}, nil
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], nil
}

// Satisfy new interface method; not used in these tests
func (f *fakePublicRepo) GetUsersBasicByIDs(ctx context.Context, ids []uint) ([]responses.UserBasic, error) {
	return []responses.UserBasic{}, nil
}

type fakeVoteRepo struct {
	hasVoted  bool
	createErr error
}

func (f *fakeVoteRepo) HasUserVoted(ctx context.Context, videoID, userID uint) (bool, error) {
	return f.hasVoted, nil
}
func (f *fakeVoteRepo) Create(ctx context.Context, videoID, userID uint) error {
	if f.createErr == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(f.createErr, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrConflict
	}
	return f.createErr
}

// Helper to setup Gin with handler and a middleware to set userID
func setupRouter(h *hdl.PublicHandlers, withAuth bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if withAuth {
		r.Use(func(c *gin.Context) { c.Set("userID", uint(10)); c.Next() })
	}
	r.POST("/api/public/videos/:video_id/vote", h.VotePublicVideo)
	r.GET("/api/public/rankings", h.ListRankings)
	return r
}

func TestVotePublicVideo_OK(t *testing.T) {
	pub := &fakePublicRepo{exists: true}
	votes := &fakeVoteRepo{hasVoted: false, createErr: nil}
	svc := useCase.NewPublicService(pub, votes)
	h := hdl.NewPublicHandlers(svc)
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
	svc := useCase.NewPublicService(pub, votes)
	h := hdl.NewPublicHandlers(svc)
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
	svc := useCase.NewPublicService(pub, votes)
	h := hdl.NewPublicHandlers(svc)
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
	svc := useCase.NewPublicService(pub, votes)
	h := hdl.NewPublicHandlers(svc)
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
	svc := useCase.NewPublicService(pub, votes)
	h := hdl.NewPublicHandlers(svc)
	r := setupRouter(h, false) // no auth middleware

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/public/videos/123/vote", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body=%s", w.Code, w.Body.String())
	}
}

// --- Rankings tests ---

func strptr(s string) *string { return &s }

func TestListRankings_OK_NoFilters(t *testing.T) {
	data := []responses.RankingItem{
		{Username: "alice", City: strptr("Bogotá"), Votes: 10},
		{Username: "bob", City: strptr("Medellín"), Votes: 5},
	}
	pub := &fakePublicRepo{exists: true, rankings: data}
	svc := useCase.NewPublicService(pub, &fakeVoteRepo{})
	h := hdl.NewPublicHandlers(svc)
	r := setupRouter(h, false)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/public/rankings", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	var got []responses.RankingEntry
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
	if got[0].Position != 1 || got[1].Position != 2 {
		t.Fatalf("expected positions 1,2 got %d,%d", got[0].Position, got[1].Position)
	}
	if got[0].Votes < got[1].Votes {
		t.Fatalf("expected sorted desc by votes, got %v", got)
	}
}

func TestListRankings_CityFilter(t *testing.T) {
	data := []responses.RankingItem{
		{Username: "alice", City: strptr("Bogotá"), Votes: 10},
		{Username: "bob", City: strptr("Medellín"), Votes: 5},
		{Username: "carl", City: strptr("Bogotá"), Votes: 3},
	}
	pub := &fakePublicRepo{exists: true, rankings: data}
	svc := useCase.NewPublicService(pub, &fakeVoteRepo{})
	h := hdl.NewPublicHandlers(svc)
	r := setupRouter(h, false)

	q := url.Values{"city": []string{"Bogotá"}}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/public/rankings?"+q.Encode(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	var got []responses.RankingEntry
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	for _, it := range got {
		if it.City == nil || *it.City != "Bogotá" {
			t.Fatalf("unexpected city in result: %+v", it)
		}
	}
}

func TestListRankings_Pagination_Page2Size1(t *testing.T) {
	data := []responses.RankingItem{
		{Username: "alice", City: strptr("Bogotá"), Votes: 10},
		{Username: "bob", City: strptr("Bogotá"), Votes: 9},
		{Username: "carl", City: strptr("Bogotá"), Votes: 8},
	}
	pub := &fakePublicRepo{exists: true, rankings: data}
	svc := useCase.NewPublicService(pub, &fakeVoteRepo{})
	h := hdl.NewPublicHandlers(svc)
	r := setupRouter(h, false)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/public/rankings?page=2&pageSize=1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	var got []responses.RankingEntry
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d (%v)", len(got), got)
	}
	if got[0].Position != 1 {
		t.Fatalf("expected position 1 in page, got %d", got[0].Position)
	}
	if got[0].Username != "bob" {
		t.Fatalf("expected second ranked username 'bob', got %s", got[0].Username)
	}
}

func TestListRankings_BadParams(t *testing.T) {
	cases := []string{
		"/api/public/rankings?page=0",
		"/api/public/rankings?page=-1",
		"/api/public/rankings?page=abc",
		"/api/public/rankings?pageSize=0",
		"/api/public/rankings?pageSize=101",
		"/api/public/rankings?pageSize=xyz",
	}
	for i, path := range cases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			pub := &fakePublicRepo{exists: true}
			svc := useCase.NewPublicService(pub, &fakeVoteRepo{})
			h := hdl.NewPublicHandlers(svc)
			r := setupRouter(h, false)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, nil)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}
