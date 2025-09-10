package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/application/useCase"
	"api/internal/domain"
	hdl "api/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
)

type fakeLocationRepo struct {
	cityID uint
	err    error
}

func (f *fakeLocationRepo) GetCityID(ctx context.Context, country, city string) (uint, error) {
	if f.err != nil {
		return 0, f.err
	}
	return f.cityID, nil
}

func setupLocationRouter(repo *fakeLocationRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := useCase.NewLocationService(repo)
	h := hdl.NewLocationHandlers(svc)
	r := gin.New()
	r.GET("/api/location/city-id", h.GetCityID)
	return r
}

func TestGetCityID_Success(t *testing.T) {
	repo := &fakeLocationRepo{cityID: 123}
	r := setupLocationRouter(repo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/location/city-id?country=Colombia&city=Bogotá", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetCityID_MissingCountry(t *testing.T) {
	repo := &fakeLocationRepo{}
	r := setupLocationRouter(repo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/location/city-id?city=Bogotá", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetCityID_MissingCity(t *testing.T) {
	repo := &fakeLocationRepo{}
	r := setupLocationRouter(repo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/location/city-id?country=Colombia", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetCityID_NotFound(t *testing.T) {
	repo := &fakeLocationRepo{err: domain.ErrNotFound}
	r := setupLocationRouter(repo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/location/city-id?country=Colombia&city=Unknown", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetCityID_Invalid(t *testing.T) {
	repo := &fakeLocationRepo{err: domain.ErrInvalid}
	r := setupLocationRouter(repo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/location/city-id?country=Invalid&city=Invalid", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}