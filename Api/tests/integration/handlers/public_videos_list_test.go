package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/application/useCase"
	"api/internal/domain/responses"
	hdl "api/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
)

type fakePublicRepoList struct {
	videos []responses.PublicVideoResponse
	err    error
}

func (f *fakePublicRepoList) ListPublicVideos(ctx context.Context) ([]responses.PublicVideoResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.videos, nil
}

func (f *fakePublicRepoList) GetPublicByID(ctx context.Context, id uint) (*responses.PublicVideoResponse, error) {
	return nil, nil
}

func (f *fakePublicRepoList) Rankings(ctx context.Context, city *string, page, pageSize int) ([]responses.RankingItem, error) {
	return nil, nil
}

func (f *fakePublicRepoList) GetUsersBasicByIDs(ctx context.Context, ids []uint) ([]responses.UserBasic, error) {
	return nil, nil
}

func setupPublicVideosRouter(repo *fakePublicRepoList) *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := useCase.NewPublicService(repo, &fakeVoteRepo{})
	h := hdl.NewPublicHandlers(svc)
	r := gin.New()
	r.GET("/api/public/videos", h.ListPublicVideos)
	return r
}

func TestListPublicVideos_Success(t *testing.T) {
	city1 := "Bogotá"
	city2 := "Medellín"
	videos := []responses.PublicVideoResponse{
		{VideoID: 1, Title: "Video 1", City: &city1, Votes: 10},
		{VideoID: 2, Title: "Video 2", City: &city2, Votes: 5},
	}
	repo := &fakePublicRepoList{videos: videos}
	r := setupPublicVideosRouter(repo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/public/videos", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	var got []responses.PublicVideoResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 videos, got %d", len(got))
	}

	if got[0].Title != "Video 1" || got[1].Title != "Video 2" {
		t.Fatalf("unexpected video titles: %+v", got)
	}
}

func TestListPublicVideos_Empty(t *testing.T) {
	repo := &fakePublicRepoList{videos: []responses.PublicVideoResponse{}}
	r := setupPublicVideosRouter(repo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/public/videos", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	var got []responses.PublicVideoResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("expected empty list, got %d videos", len(got))
	}
}
