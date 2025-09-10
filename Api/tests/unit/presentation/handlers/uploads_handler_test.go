package handlers_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api/internal/application/useCase"
	"api/internal/domain/entities"
	"api/internal/domain/requests"
	"api/internal/domain/responses"
	"api/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockUploadsRepo struct {
	err error
}

func (m *mockUploadsRepo) Create(ctx context.Context, video *entities.Video) error { return m.err }
func (m *mockUploadsRepo) GetByID(ctx context.Context, id uint) (*entities.Video, error) { return nil, nil }
func (m *mockUploadsRepo) List(ctx context.Context) ([]*entities.Video, error) { return nil, nil }
func (m *mockUploadsRepo) ListByUser(ctx context.Context, userID uint) ([]*entities.Video, error) { return nil, nil }
func (m *mockUploadsRepo) GetByIDAndUser(ctx context.Context, id, userID uint) (*entities.Video, error) { return nil, nil }
func (m *mockUploadsRepo) Delete(ctx context.Context, id uint) error { return nil }
func (m *mockUploadsRepo) UpdateStatus(ctx context.Context, id uint, status entities.VideoStatus) error { return nil }

type mockUploadsStorage struct {
	policy *responses.CreateUploadResponsePostPolicy
	err    error
}

func (m *mockUploadsStorage) Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	return "", nil
}

func (m *mockUploadsStorage) PresignedPostPolicy(ctx context.Context, req requests.CreateUploadRequest) (*responses.CreateUploadResponsePostPolicy, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.policy, nil
}

func TestUploadsHandler_CreatePostPolicy_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	policy := &responses.CreateUploadResponsePostPolicy{
		URL: "https://example.com/upload",
		Fields: map[string]string{
			"key":    "test-key",
			"policy": "test-policy",
		},
	}
	repo := &mockUploadsRepo{}
	storage := &mockUploadsStorage{policy: policy}
	uc := useCase.NewUploadsUseCase(repo, storage, nil, "")
	h := handlers.NewUploadsHandlers(uc)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(10))
		c.Next()
	})
	r.POST("/api/uploads", h.CreatePostPolicy)

	body := `{"title":"Test Video","status":"UPLOADED"}`
	req := httptest.NewRequest(http.MethodPost, "/api/uploads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "https://example.com/upload")
}

func TestUploadsHandler_CreatePostPolicy_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := &mockUploadsRepo{}
	storage := &mockUploadsStorage{}
	uc := useCase.NewUploadsUseCase(repo, storage, nil, "")
	h := handlers.NewUploadsHandlers(uc)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(10))
		c.Next()
	})
	r.POST("/api/uploads", h.CreatePostPolicy)

	body := `{"title":"Test Video"`
	req := httptest.NewRequest(http.MethodPost, "/api/uploads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadsHandler_CreatePostPolicy_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := &mockUploadsRepo{}
	storage := &mockUploadsStorage{}
	uc := useCase.NewUploadsUseCase(repo, storage, nil, "")
	h := handlers.NewUploadsHandlers(uc)

	r.POST("/api/uploads", h.CreatePostPolicy)

	body := `{"title":"Test Video","status":"UPLOADED"}`
	req := httptest.NewRequest(http.MethodPost, "/api/uploads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUploadsHandler_CreatePostPolicy_StorageError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := &mockUploadsRepo{}
	storage := &mockUploadsStorage{err: errors.New("storage error")}
	uc := useCase.NewUploadsUseCase(repo, storage, nil, "")
	h := handlers.NewUploadsHandlers(uc)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(10))
		c.Next()
	})
	r.POST("/api/uploads", h.CreatePostPolicy)

	body := `{"title":"Test Video","status":"UPLOADED"}`
	req := httptest.NewRequest(http.MethodPost, "/api/uploads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUploadsHandler_CreatePostPolicy_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	policy := &responses.CreateUploadResponsePostPolicy{
		URL: "https://example.com/upload",
		Fields: map[string]string{
			"key":    "test-key",
			"policy": "test-policy",
		},
	}
	repo := &mockUploadsRepo{err: errors.New("database error")}
	storage := &mockUploadsStorage{policy: policy}
	uc := useCase.NewUploadsUseCase(repo, storage, nil, "")
	h := handlers.NewUploadsHandlers(uc)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(10))
		c.Next()
	})
	r.POST("/api/uploads", h.CreatePostPolicy)

	body := `{"title":"Test Video","status":"UPLOADED"}`
	req := httptest.NewRequest(http.MethodPost, "/api/uploads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}