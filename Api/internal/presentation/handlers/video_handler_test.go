package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"api/internal/application/useCase"
	"api/internal/domain/entities"
	"api/internal/domain/requests"
	"api/internal/domain/responses"

	"github.com/gin-gonic/gin"
)

// --- Fakes ---

type fakeVideoRepo struct {
	list []*entities.Video
	err  error
}

func (f *fakeVideoRepo) Create(_ context.Context, _ *entities.Video) error          { return nil }
func (f *fakeVideoRepo) GetByID(_ context.Context, _ uint) (*entities.Video, error) { return nil, nil }
func (f *fakeVideoRepo) List(_ context.Context) ([]*entities.Video, error)          { return nil, nil }
func (f *fakeVideoRepo) ListByUser(_ context.Context, _ uint) ([]*entities.Video, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.list, nil
}

type fakeStorage struct{}

func (f *fakeStorage) Save(_ context.Context, _ string, _ io.Reader, _ int64, _ string) (string, error) {
	return "", nil
}
func (f *fakeStorage) PresignedPostPolicy(_ context.Context, _ requests.CreateUploadRequest) (*responses.CreateUploadResponsePostPolicy, error) {
	return nil, nil
}

func setupVideoRouter(uc *useCase.UploadsUseCase, withAuth bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewVideoHandlers(uc)
	r := gin.New()
	if withAuth {
		r.Use(func(c *gin.Context) { c.Set("userID", uint(10)); c.Next() })
	}
	r.GET("/api/videos", h.ListVideos)
	return r
}

func TestListVideos_OK(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	processed := now.Add(5 * time.Minute)
	processedURL := "https://cdn/x.mp4"
	repo := &fakeVideoRepo{list: []*entities.Video{
		{VideoID: 1, UserID: 10, Title: "A", OriginalFile: "s3://orig/a.mp4", ProcessedFile: &processedURL, Status: string(entities.StatusProcessed), UploadedAt: now, ProcessedAt: &processed},
		{VideoID: 2, UserID: 10, Title: "B", OriginalFile: "s3://orig/b.mp4", Status: string(entities.StatusUploaded), UploadedAt: now.Add(1 * time.Minute)},
	}}
	uc := useCase.NewUploadsUseCase(repo, &fakeStorage{})
	r := setupVideoRouter(uc, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/videos", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", w.Code, w.Body.String())
	}
	var got []responses.VideoResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
	if got[0].VideoID != "1" || got[0].Status != "processed" || got[0].ProcessedURL == nil {
		t.Fatalf("unexpected first item: %+v", got[0])
	}
	if got[1].VideoID != "2" || got[1].Status != "uploaded" || got[1].ProcessedURL != nil {
		t.Fatalf("unexpected second item: %+v", got[1])
	}
}

func TestListVideos_Empty(t *testing.T) {
	repo := &fakeVideoRepo{list: []*entities.Video{}}
	uc := useCase.NewUploadsUseCase(repo, &fakeStorage{})
	r := setupVideoRouter(uc, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/videos", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", w.Code, w.Body.String())
	}
	var got []responses.VideoResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty list, got %d", len(got))
	}
}

func TestListVideos_RepoError(t *testing.T) {
	repo := &fakeVideoRepo{err: errors.New("boom")}
	uc := useCase.NewUploadsUseCase(repo, &fakeStorage{})
	r := setupVideoRouter(uc, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/videos", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d (%s)", w.Code, w.Body.String())
	}
}

func TestListVideos_Unauthorized(t *testing.T) {
	repo := &fakeVideoRepo{}
	uc := useCase.NewUploadsUseCase(repo, &fakeStorage{})
	r := setupVideoRouter(uc, false) // no auth middleware sets userID

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/videos", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d (%s)", w.Code, w.Body.String())
	}
}
