package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockLocationRepo struct {
	cityID uint
	err    error
}

func (m *mockLocationRepo) GetCityID(ctx context.Context, country, city string) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	return int(m.cityID), nil
}

func TestLocationHandlers_GetCityID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := &mockLocationRepo{cityID: 42}
	svc := useCase.NewLocationService(repo)
	h := handlers.NewLocationHandlers(svc)

	r.GET("/api/location/city-id", h.GetCityID)

	req := httptest.NewRequest(http.MethodGet, "/api/location/city-id?country=Colombia&city=Bogotá", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"city_id\":42")
}

func TestLocationHandlers_GetCityID_MissingParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := &mockLocationRepo{}
	svc := useCase.NewLocationService(repo)
	h := handlers.NewLocationHandlers(svc)

	r.GET("/api/location/city-id", h.GetCityID)

	tests := []struct {
		name string
		url  string
	}{
		{"missing country", "/api/location/city-id?city=Bogotá"},
		{"missing city", "/api/location/city-id?country=Colombia"},
		{"empty country", "/api/location/city-id?country=&city=Bogotá"},
		{"empty city", "/api/location/city-id?country=Colombia&city="},
		{"whitespace only", "/api/location/city-id?country=%20%20%20&city=%20%20%20"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "country and city are required")
		})
	}
}

func TestLocationHandlers_GetCityID_ErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedMsg    string
	}{
		{"invalid input", domain.ErrInvalid, http.StatusBadRequest, "invalid country or city"},
		{"not found", domain.ErrNotFound, http.StatusNotFound, "city not found for country"},
		{"internal error", errors.New("db error"), http.StatusInternalServerError, "db error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			repo := &mockLocationRepo{err: tt.err}
			svc := useCase.NewLocationService(repo)
			h := handlers.NewLocationHandlers(svc)

			r.GET("/api/location/city-id", h.GetCityID)

			req := httptest.NewRequest(http.MethodGet, "/api/location/city-id?country=Colombia&city=Bogotá", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedMsg)
		})
	}
}