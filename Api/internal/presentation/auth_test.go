package presentation

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"main_videork/internal/application/useCase"
	"main_videork/internal/domain"
	"main_videork/internal/domain/entities"
)

type mockUserRepo struct{}

func (m *mockUserRepo) Create(ctx context.Context, user *entities.User) error {
	return nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if email == "exists@example.com" {
		return &entities.User{}, nil
	}
	return nil, domain.ErrNotFound
}

func TestAuthHandlers_RegisterEmailExists(t *testing.T) {
	repo := &mockUserRepo{}
	service := useCase.NewAuthService(repo, "secret")
	handler := NewAuthHandlers(service)

	router := gin.New()
	router.POST("/api/auth/signup", handler.Register)

	body := []byte(`{"first_name":"A","last_name":"B","email":"exists@example.com","password1":"p","password2":"p"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}
