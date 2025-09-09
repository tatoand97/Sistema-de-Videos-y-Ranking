package handlers_test

import (
	"api/internal/application/useCase"
	"api/internal/presentation/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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
