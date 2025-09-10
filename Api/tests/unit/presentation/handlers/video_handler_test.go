package handlers_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/entities"
	"api/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockVideoRepo struct {
	videos []*entities.Video
	video  *entities.Video
	err    error
}

func (m *mockVideoRepo) Create(ctx context.Context, video *entities.Video) error { return nil }
func (m *mockVideoRepo) GetByID(ctx context.Context, id uint) (*entities.Video, error) {
	return m.video, m.err
}
func (m *mockVideoRepo) List(ctx context.Context) ([]*entities.Video, error) { return nil, nil }
func (m *mockVideoRepo) Delete(ctx context.Context, id uint) error           { return m.err }
func (m *mockVideoRepo) UpdateStatus(ctx context.Context, id uint, status entities.VideoStatus) error {
	return m.err
}

func (m *mockVideoRepo) ListByUser(ctx context.Context, userID uint) ([]*entities.Video, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.videos, nil
}

func (m *mockVideoRepo) GetByIDAndUser(ctx context.Context, id, userID uint) (*entities.Video, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.video, nil
}

type mockVideoStorage struct {
	err error
}

func (m *mockVideoStorage) Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	return "", nil
}

func TestVideoHandlers_ListVideos_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	videos := []*entities.Video{
		{VideoID: 1, UserID: 10, Title: "Video 1", Status: "UPLOADED"},
		{VideoID: 2, UserID: 10, Title: "Video 2", Status: "PROCESSED"},
	}
	repo := &mockVideoRepo{videos: videos}
	storage := &mockVideoStorage{}
	uc := useCase.NewUploadsUseCase(repo, storage, nil, "")
	h := handlers.NewVideoHandlers(uc)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(10))
		c.Next()
	})
	r.GET("/api/videos", h.ListVideos)

	req := httptest.NewRequest(http.MethodGet, "/api/videos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Video 1")
	assert.Contains(t, w.Body.String(), "Video 2")
}

func TestVideoHandlers_ListVideos_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := &mockVideoRepo{}
	storage := &mockVideoStorage{}
	uc := useCase.NewUploadsUseCase(repo, storage, nil, "")
	h := handlers.NewVideoHandlers(uc)

	r.GET("/api/videos", h.ListVideos)

	req := httptest.NewRequest(http.MethodGet, "/api/videos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestVideoHandlers_GetVideoDetail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	video := &entities.Video{
		VideoID: 123,
		UserID:  10,
		Title:   "Test Video",
		Status:  "PROCESSED",
	}
	repo := &mockVideoRepo{video: video}
	storage := &mockVideoStorage{}
	uc := useCase.NewUploadsUseCase(repo, storage, nil, "")
	h := handlers.NewVideoHandlers(uc)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(10))
		c.Next()
	})
	r.GET("/api/videos/:video_id", h.GetVideoDetail)

	req := httptest.NewRequest(http.MethodGet, "/api/videos/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Video")
}

func TestVideoHandlers_GetVideoDetail_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := &mockVideoRepo{err: domain.ErrNotFound}
	storage := &mockVideoStorage{}
	uc := useCase.NewUploadsUseCase(repo, storage, nil, "")
	h := handlers.NewVideoHandlers(uc)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(10))
		c.Next()
	})
	r.GET("/api/videos/:video_id", h.GetVideoDetail)

	req := httptest.NewRequest(http.MethodGet, "/api/videos/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestVideoHandlers_DeleteVideo_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	video := &entities.Video{VideoID: 123, UserID: 10}
	repo := &mockVideoRepo{video: video}
	storage := &mockVideoStorage{}
	uc := useCase.NewUploadsUseCase(repo, storage, nil, "")
	h := handlers.NewVideoHandlers(uc)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(10))
		c.Next()
	})
	r.DELETE("/api/videos/:video_id", h.DeleteVideo)

	req := httptest.NewRequest(http.MethodDelete, "/api/videos/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "eliminado")
}

func TestVideoHandlers_PublishVideo_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	processed := "https://example.com/processed.mp4"
	video := &entities.Video{
		VideoID:       123,
		UserID:        10,
		Status:        string(entities.StatusProcessed),
		ProcessedFile: &processed,
	}
	repo := &mockVideoRepo{video: video}
	storage := &mockVideoStorage{}
	uc := useCase.NewUploadsUseCase(repo, storage, nil, "")
	h := handlers.NewVideoHandlers(uc)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(10))
		c.Set("permissions", []string{"edit_video"})
		c.Next()
	})
	r.POST("/api/videos/:video_id/publish", h.PublishVideo)

	req := httptest.NewRequest(http.MethodPost, "/api/videos/123/publish", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVideoHandlers_PublishVideo_NotProcessed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	video := &entities.Video{
		VideoID: 123,
		UserID:  10,
		Status:  string(entities.StatusUploaded),
	}
	repo := &mockVideoRepo{video: video}
	storage := &mockVideoStorage{}
	uc := useCase.NewUploadsUseCase(repo, storage, nil, "")
	h := handlers.NewVideoHandlers(uc)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(10))
		c.Set("permissions", []string{"moderate_content"})
		c.Next()
	})
	r.POST("/api/videos/:video_id/publish", h.PublishVideo)

	req := httptest.NewRequest(http.MethodPost, "/api/videos/123/publish", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
