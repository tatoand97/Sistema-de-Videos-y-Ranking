package handlers_test

import (
	"api/internal/application/useCase"
	"api/internal/presentation/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStatusHandlers_ListVideoStatuses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := useCase.NewStatusService()
	h := handlers.NewStatusHandlers(svc)
	r.GET("/api/videos/statuses", h.ListVideoStatuses)

	req := httptest.NewRequest(http.MethodGet, "/api/videos/statuses", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string][]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	expected := []string{"UPLOADED", "TRIMMING", "ADJUSTING_RESOLUTION", "ADDING_WATERMARK", "REMOVING_AUDIO", "ADDING_INTRO_OUTRO", "PROCESSED", "FAILED"}
	assert.Equal(t, expected, body["statuses"])
}
