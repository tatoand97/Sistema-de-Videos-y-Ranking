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
	"api/internal/domain/entities"
	"api/internal/domain/responses"
	hdl "api/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
)

type fakeVideoRepoDetail struct {
	video *entities.Video
	err   error
}

func (f *fakeVideoRepoDetail) Create(_ context.Context, _ *entities.Video) error { return nil }
func (f *fakeVideoRepoDetail) GetByID(_ context.Context, _ uint) (*entities.Video, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.video, nil
}
func (f *fakeVideoRepoDetail) List(_ context.Context) ([]*entities.Video, error) { return nil, nil }
func (f *fakeVideoRepoDetail) ListByUser(_ context.Context, _ uint) ([]*entities.Video, error) {
	return nil, nil
}
func (f *fakeVideoRepoDetail) Delete(_ context.Context, _ uint) error { return nil }
func (f *fakeVideoRepoDetail) UpdateStatus(_ context.Context, _ uint, _ entities.VideoStatus) error {
	return nil
}

func (f *fakeVideoRepoDetail) GetByIDAndUser(_ context.Context, id, userID uint) (*entities.Video, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.video, nil
}

func setupVideoDetailRouter(repo *fakeVideoRepoDetail, withAuth bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	uc := useCase.NewUploadsUseCase(repo, &fakeStorage{}, nil, "")
	h := hdl.NewVideoHandlers(uc, processedBaseURL)
	r := gin.New()
	if withAuth {
		r.Use(func(c *gin.Context) {
			c.Set("userID", uint(10))
			c.Set("permissions", []string{"edit_video"})
			c.Next()
		})
	}
	r.GET("/api/videos/:video_id", h.GetVideoDetail)
	r.DELETE("/api/videos/:video_id", h.DeleteVideo)
	r.POST("/api/videos/:video_id/publish", h.PublishVideo)
	return r
}

func TestGetVideoDetail_Success(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	processedURL := "https://cdn/processed.mp4"
	video := &entities.Video{
		VideoID:       123,
		UserID:        10,
		Title:         "Test Video",
		OriginalFile:  "s3://orig/test.mp4",
		ProcessedFile: &processedURL,
		Status:        string(entities.StatusProcessed),
		UploadedAt:    now,
		ProcessedAt:   &now,
	}
	repo := &fakeVideoRepoDetail{video: video}
	r := setupVideoDetailRouter(repo, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/videos/123", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	var got responses.VideoResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got.VideoID != "123" || got.Title != "Test Video" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestGetVideoDetail_NotFound(t *testing.T) {
	repo := &fakeVideoRepoDetail{err: domain.ErrNotFound}
	r := setupVideoDetailRouter(repo, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/videos/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetVideoDetail_Forbidden(t *testing.T) {
	repo := &fakeVideoRepoDetail{err: domain.ErrForbidden}
	r := setupVideoDetailRouter(repo, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/videos/123", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetVideoDetail_Unauthorized(t *testing.T) {
	repo := &fakeVideoRepoDetail{}
	r := setupVideoDetailRouter(repo, false)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/videos/123", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDeleteVideo_Success(t *testing.T) {
	video := &entities.Video{VideoID: 123, UserID: 10}
	repo := &fakeVideoRepoDetail{video: video}
	r := setupVideoDetailRouter(repo, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/videos/123", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDeleteVideo_NotFound(t *testing.T) {
	repo := &fakeVideoRepoDetail{err: domain.ErrNotFound}
	r := setupVideoDetailRouter(repo, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/videos/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestPublishVideo_Success(t *testing.T) {
	processed := "https://cdn/processed.mp4"
	video := &entities.Video{
		VideoID:       123,
		UserID:        10,
		Status:        string(entities.StatusProcessed),
		ProcessedFile: &processed,
	}
	repo := &fakeVideoRepoDetail{video: video}
	r := setupVideoDetailRouter(repo, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/videos/123/publish", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestPublishVideo_NotProcessed(t *testing.T) {
	video := &entities.Video{
		VideoID: 123,
		UserID:  10,
		Status:  string(entities.StatusUploaded),
	}
	repo := &fakeVideoRepoDetail{video: video}
	r := setupVideoDetailRouter(repo, true)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/videos/123/publish", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}
