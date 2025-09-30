package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/responses"
	"api/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestPublicHandlers_VotePublicVideo_MissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := useCase.NewPublicService(nil, nil)
	h := handlers.NewPublicHandlers(svc)
	r.POST("/api/public/videos/:video_id/vote", h.VotePublicVideo)

	req := httptest.NewRequest(http.MethodPost, "/api/public/videos/123/vote", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPublicHandlers_ListRankings_InvalidParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := useCase.NewPublicService(nil, nil)
	h := handlers.NewPublicHandlers(svc)
	r.GET("/api/public/rankings", h.ListRankings)

	req := httptest.NewRequest(http.MethodGet, "/api/public/rankings?page=zero", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPublicHandlers_ListRankings_UsesCache(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &cacheRepo{
		users: map[uint]responses.UserBasic{
			1: {UserID: 1, Username: "alice", City: strPtr("Bogotá")},
			2: {UserID: 2, Username: "bob", City: strPtr("Medellín")},
		},
	}
	svc := useCase.NewPublicService(repo, nil)
	cache := &fakeCache{data: make(map[string][]byte)}

	now := time.Now().UTC()
	entry := map[string]any{
		"schema_version": "v2",
		"scope":          "global",
		"as_of":          now.Format(time.RFC3339),
		"fresh_until":    now.Add(2 * time.Minute).Format(time.RFC3339),
		"stale_until":    now.Add(10 * time.Minute).Format(time.RFC3339),
		"items": []map[string]any{
			{"rank": 1, "user_id": 1, "username": "alice", "score": 12},
			{"rank": 2, "user_id": 2, "username": "bob", "score": 9},
		},
	}
	payload, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("failed to marshal cache entry: %v", err)
	}
	cache.data["rank:global:v2"] = payload

	h := handlers.NewPublicHandlersWithCache(svc, cache, "v2")
	r := gin.New()
	r.GET("/api/public/rankings", h.ListRankings)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/public/rankings", nil)
	r.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusOK, w.Code) {
		return
	}

	var got []responses.RankingEntry
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if assert.Len(t, got, 2) {
		assert.Equal(t, "alice", got[0].Username)
		if assert.NotNil(t, got[0].City) {
			assert.Equal(t, "Bogotá", *got[0].City)
		}
		assert.Equal(t, 0, repo.rankingsCalls)
	}
}

func TestPublicHandlers_ListRankings_CacheCityScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &cacheRepo{
		users: map[uint]responses.UserBasic{
			9: {UserID: 9, Username: "caro", City: strPtr("Bogotá")},
		},
	}
	svc := useCase.NewPublicService(repo, nil)
	cache := &fakeCache{data: make(map[string][]byte)}

	now := time.Now().UTC()
	entry := map[string]any{
		"schema_version": "v2",
		"scope":          "city",
		"city":           "Bogotá",
		"city_slug":      "bogota",
		"as_of":          now.Format(time.RFC3339),
		"fresh_until":    now.Add(2 * time.Minute).Format(time.RFC3339),
		"stale_until":    now.Add(10 * time.Minute).Format(time.RFC3339),
		"items": []map[string]any{
			{"rank": 1, "user_id": 9, "username": "caro", "score": 25},
		},
	}
	payload, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("failed to marshal cache entry: %v", err)
	}
	cache.data["rank:city:bogota:v2"] = payload

	h := handlers.NewPublicHandlersWithCache(svc, cache, "v2")
	r := gin.New()
	r.GET("/api/public/rankings", h.ListRankings)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/public/rankings?city=Bogotá", nil)
	r.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusOK, w.Code) {
		return
	}

	var got []responses.RankingEntry
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if assert.Len(t, got, 1) {
		assert.Equal(t, "caro", got[0].Username)
		if assert.NotNil(t, got[0].City) {
			assert.Equal(t, "Bogotá", *got[0].City)
		}
		assert.Equal(t, 0, repo.rankingsCalls)
	}
}

func TestPublicHandlers_ListRankings_CacheStaleFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &cacheRepo{
		users:           map[uint]responses.UserBasic{},
		rankingResponse: []responses.RankingItem{{Username: "db", Votes: 3}},
	}
	svc := useCase.NewPublicService(repo, nil)
	cache := &fakeCache{data: make(map[string][]byte)}

	now := time.Now().UTC()
	entry := map[string]any{
		"schema_version": "v2",
		"scope":          "global",
		"as_of":          now.Add(-10 * time.Minute).Format(time.RFC3339),
		"fresh_until":    now.Add(-5 * time.Minute).Format(time.RFC3339),
		"stale_until":    now.Add(-1 * time.Minute).Format(time.RFC3339),
		"items": []map[string]any{
			{"rank": 1, "user_id": 1, "username": "alice", "score": 12},
		},
	}
	payload, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("failed to marshal cache entry: %v", err)
	}
	cache.data["rank:global:v2"] = payload

	h := handlers.NewPublicHandlersWithCache(svc, cache, "v2")
	r := gin.New()
	r.GET("/api/public/rankings", h.ListRankings)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/public/rankings", nil)
	r.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusOK, w.Code) {
		return
	}

	var got []responses.RankingEntry
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if assert.Len(t, got, 1) {
		assert.Equal(t, "db", got[0].Username)
	}
	assert.Equal(t, 1, repo.rankingsCalls)
}

func strPtr(v string) *string {
	out := v
	return &out
}

type fakeCache struct {
	data map[string][]byte
}

func (f *fakeCache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	if f.data == nil {
		return nil, redis.Nil
	}
	if v, ok := f.data[key]; ok {
		return v, nil
	}
	return nil, redis.Nil
}

type cacheRepo struct {
	users           map[uint]responses.UserBasic
	rankingResponse []responses.RankingItem
	rankingsCalls   int
}

func (r *cacheRepo) ListPublicVideos(ctx context.Context) ([]responses.PublicVideoResponse, error) {
	return []responses.PublicVideoResponse{}, nil
}

func (r *cacheRepo) GetPublicByID(ctx context.Context, id uint) (*responses.PublicVideoResponse, error) {
	return nil, domain.ErrNotFound
}

func (r *cacheRepo) Rankings(ctx context.Context, city *string, page, pageSize int) ([]responses.RankingItem, error) {
	r.rankingsCalls++
	if len(r.rankingResponse) == 0 {
		return []responses.RankingItem{}, nil
	}
	return r.rankingResponse, nil
}

func (r *cacheRepo) GetUsersBasicByIDs(ctx context.Context, ids []uint) ([]responses.UserBasic, error) {
	if len(ids) == 0 {
		return []responses.UserBasic{}, nil
	}
	res := make([]responses.UserBasic, 0, len(ids))
	for _, id := range ids {
		if ub, ok := r.users[id]; ok {
			res = append(res, ub)
		}
	}
	return res, nil
}
