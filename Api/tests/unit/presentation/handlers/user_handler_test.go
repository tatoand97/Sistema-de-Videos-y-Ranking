package handlers_test

import (
	usecase "api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/entities"
	"api/internal/presentation/handlers"
	"api/tests/mocks"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserHandlers_Register_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	userRepo := &mocks.MockUserRepository{}
	svc := usecase.NewUserService(userRepo, nil)
	h := handlers.NewUserHandlers(svc)
	r.POST("/api/auth/signup", h.Register)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandlers_Register_EmailInUse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	userRepo := &mocks.MockUserRepository{GetByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
		return &entities.User{Email: email}, nil
	}}
	svc := usecase.NewUserService(userRepo, nil)
	h := handlers.NewUserHandlers(svc)
	r.POST("/api/auth/signup", h.Register)

	body := map[string]any{
		"first_name": "Ana",
		"last_name":  "Lopez",
		"email":      "ana@example.com",
		"password1":  "abc123",
		"password2":  "abc123",
		"country":    "Peru",
		"city":       "Lima",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandlers_Register_PasswordsDoNotMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	userRepo := &mocks.MockUserRepository{GetByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
		return nil, domain.ErrNotFound
	}}
	svc := usecase.NewUserService(userRepo, nil)
	h := handlers.NewUserHandlers(svc)
	r.POST("/api/auth/signup", h.Register)

	body := map[string]any{
		"first_name": "Ana",
		"last_name":  "Lopez",
		"email":      "ana@example.com",
		"password1":  "abc123",
		"password2":  "xyz789",
		"country":    "Peru",
		"city":       "Lima",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
