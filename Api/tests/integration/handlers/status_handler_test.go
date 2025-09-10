package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/application/useCase"
	hdl "api/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
)

func TestStatusHandler_ListVideoStatuses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := useCase.NewStatusService()
	h := hdl.NewStatusHandlers(svc)
	r := gin.New()
	r.GET("/api/videos/statuses", h.ListVideoStatuses)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/videos/statuses", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	var response map[string][]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	statuses, exists := response["statuses"]
	if !exists {
		t.Fatal("expected 'statuses' field in response")
	}

	expectedStatuses := []string{
		"UPLOADED", "TRIMMING", "ADJUSTING_RESOLUTION", 
		"ADDING_WATERMARK", "REMOVING_AUDIO", "ADDING_INTRO_OUTRO", 
		"PROCESSED", "PUBLISHED", "FAILED",
	}

	if len(statuses) != len(expectedStatuses) {
		t.Fatalf("expected %d statuses, got %d", len(expectedStatuses), len(statuses))
	}

	for i, expected := range expectedStatuses {
		if statuses[i] != expected {
			t.Fatalf("expected status %s at position %d, got %s", expected, i, statuses[i])
		}
	}
}