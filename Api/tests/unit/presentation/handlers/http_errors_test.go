package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/application/useCase"
	"api/internal/domain"
	hdl "api/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// These tests assert error-to-status mapping using existing handlers
// (e.g., LocationHandlers) instead of a removed HandleError helper.

type fakeLocRepoErr struct{ err error }

func (f *fakeLocRepoErr) GetCityID(ctx context.Context, country, city string) (int, error) {
	return 0, f.err
}

func TestErrorMapping_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := useCase.NewLocationService(&fakeLocRepoErr{err: domain.ErrNotFound})
	h := hdl.NewLocationHandlers(svc)
	r.GET("/city", h.GetCityID)

	req := httptest.NewRequest(http.MethodGet, "/city?country=CO&city=XX", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Not Found")
}

func TestErrorMapping_Invalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := useCase.NewLocationService(&fakeLocRepoErr{err: domain.ErrInvalid})
	h := hdl.NewLocationHandlers(svc)
	r.GET("/city", h.GetCityID)

	req := httptest.NewRequest(http.MethodGet, "/city?country=CO&city=XX", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestErrorMapping_Internal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := useCase.NewLocationService(&fakeLocRepoErr{err: errors.New("db")})
	h := hdl.NewLocationHandlers(svc)
	r.GET("/city", h.GetCityID)

	req := httptest.NewRequest(http.MethodGet, "/city?country=CO&city=XX", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
