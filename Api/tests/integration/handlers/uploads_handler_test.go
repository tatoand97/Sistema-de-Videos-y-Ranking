package handlers_test

import (
	"api/internal/application/useCase"
	"api/internal/presentation/handlers"
	"api/tests/mocks"
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUploadsHandler_UploadVideo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMocks     func() (*mocks.MockVideoRepository, *mocks.MockVideoStorage, *mocks.MockMessagePublisher)
		setupRequest   func() *http.Request
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name: "successful upload",
			setupMocks: func() (*mocks.MockVideoRepository, *mocks.MockVideoStorage, *mocks.MockMessagePublisher) {
				repo := &mocks.MockVideoRepository{
					CreateFunc: func(ctx context.Context, video interface{}) error {
						return nil
					},
				}
				storage := &mocks.MockVideoStorage{
					SaveFunc: func(ctx context.Context, objectName string, reader interface{}, size int64, contentType string) (string, error) {
						return "https://example.com/video.mp4", nil
					},
				}
				publisher := &mocks.MockMessagePublisher{}
				return repo, storage, publisher
			},
			setupRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.WriteField("title", "Test Video")
				writer.WriteField("status", "UPLOADED")
				
				part, _ := writer.CreateFormFile("file", "test.mp4")
				part.Write(createValidMP4())
				writer.Close()

				req := httptest.NewRequest("POST", "/api/uploads", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body string) {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "video_id")
				assert.Contains(t, response, "title")
			},
		},
		{
			name: "missing file",
			setupMocks: func() (*mocks.MockVideoRepository, *mocks.MockVideoStorage, *mocks.MockMessagePublisher) {
				return &mocks.MockVideoRepository{}, &mocks.MockVideoStorage{}, &mocks.MockMessagePublisher{}
			},
			setupRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.WriteField("title", "Test Video")
				writer.Close()

				req := httptest.NewRequest("POST", "/api/uploads", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, storage, publisher := tt.setupMocks()
			usecase := useCase.NewUploadsUseCase(repo, storage, publisher, "test-queue")
			handler := handlers.NewUploadsHandler(usecase)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("userID", uint(1))
				c.Next()
			})
			router.POST("/api/uploads", handler.UploadVideo)

			req := tt.setupRequest()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

func createValidMP4() []byte {
	return []byte{
		0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70,
		0x69, 0x73, 0x6f, 0x6d, 0x00, 0x00, 0x02, 0x00,
		0x69, 0x73, 0x6f, 0x6d, 0x69, 0x73, 0x6f, 0x32,
		0x61, 0x76, 0x63, 0x31, 0x6d, 0x70, 0x34, 0x31,
	}
}