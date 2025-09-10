package handlers_test

import (
	"api/internal/application/useCase"
	"api/internal/presentation/handlers"
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api/internal/domain/entities"
	"api/internal/domain/responses"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockUploadsRepo struct {
	err error
}

func (m *mockUploadsRepo) Create(ctx context.Context, video *entities.Video) error { return m.err }
func (m *mockUploadsRepo) GetByID(ctx context.Context, id uint) (*entities.Video, error) {
	return nil, nil
}
func (m *mockUploadsRepo) List(ctx context.Context) ([]*entities.Video, error) { return nil, nil }
func (m *mockUploadsRepo) ListByUser(ctx context.Context, userID uint) ([]*entities.Video, error) {
	return nil, nil
}
func (m *mockUploadsRepo) GetByIDAndUser(ctx context.Context, id, userID uint) (*entities.Video, error) {
	return nil, nil
}
func (m *mockUploadsRepo) Delete(ctx context.Context, id uint) error { return nil }
func (m *mockUploadsRepo) UpdateStatus(ctx context.Context, id uint, status entities.VideoStatus) error {
	return nil
}

type mockUploadsStorage struct {
	policy *responses.CreateUploadResponsePostPolicy
	err    error
}

func (m *mockUploadsStorage) Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	return "https://example.com/video.mp4", nil
}

func multipartWithFile(fields map[string]string, filename string, content []byte) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range fields {
		writer.WriteField(k, v)
	}
	part, _ := writer.CreateFormFile("file", filename)
	part.Write(content)
	writer.Close()
	return body, writer.FormDataContentType()
}

func TestUploadsHandler_CreatePostPolicy_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	policy := &responses.CreateUploadResponsePostPolicy{
		UploadURL:   "https://example.com/upload",
		ResourceURL: "https://example.com/resource",
		ExpiresAt:   "2024-01-01T01:00:00Z",
		Form: responses.S3PostPolicyForm{
			Key:    "test-key",
			Policy: "test-policy",
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

	body := `{"filename":"test.mp4","mimeType":"video/mp4","sizeBytes":1024,"checksum":"abc123"}`
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

	body := `{"filename":"test.mp4"`
	req := httptest.NewRequest(http.MethodPost, "/api/uploads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code) // invalid mp4 will trigger Bad Request via validations
}

func TestUploadsHandler_UploadVideo_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	uc := useCase.NewUploadsUseCase(&mockUploadsRepo{}, &mockUploadsStorage{}, nil, "")
	h := handlers.NewUploadsHandlers(uc)

	r.POST("/api/uploads", h.CreatePostPolicy)

	body := `{"filename":"test.mp4","mimeType":"video/mp4","sizeBytes":1024,"checksum":"abc123"}`
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

	body := `{"filename":"test.mp4","mimeType":"video/mp4","sizeBytes":1024,"checksum":"abc123"}`
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
		UploadURL:   "https://example.com/upload",
		ResourceURL: "https://example.com/resource",
		ExpiresAt:   "2024-01-01T01:00:00Z",
		Form: responses.S3PostPolicyForm{
			Key:    "test-key",
			Policy: "test-policy",
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

	body := `{"filename":"test.mp4","mimeType":"video/mp4","sizeBytes":1024,"checksum":"abc123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/uploads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}