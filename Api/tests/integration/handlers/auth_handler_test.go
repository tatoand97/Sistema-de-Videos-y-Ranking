package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/entities"
	hdl "api/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepoAuth struct {
	user *entities.User
	err  error
}

func (f *fakeUserRepoAuth) Create(ctx context.Context, user *entities.User) error { return nil }
func (f *fakeUserRepoAuth) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.user, nil
}
func (f *fakeUserRepoAuth) EmailExists(ctx context.Context, email string) (bool, error) {
	return false, nil
}
func (f *fakeUserRepoAuth) GetPermissions(ctx context.Context, userID uint) ([]string, error) {
	return []string{"read"}, nil
}

type fakeCacheAuth struct {
	blacklisted bool
}

func (f *fakeCacheAuth) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	return nil
}
func (f *fakeCacheAuth) Get(ctx context.Context, key string) (string, error) { return "", nil }
func (f *fakeCacheAuth) Delete(ctx context.Context, key string) error        { return nil }
func (f *fakeCacheAuth) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	return f.blacklisted, nil
}
func (f *fakeCacheAuth) BlacklistToken(ctx context.Context, token string, ttl int) error { return nil }

func setupAuthRouter(userRepo *fakeUserRepoAuth, cache *fakeCacheAuth) *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := useCase.NewAuthService(userRepo, "secret")
	h := hdl.NewAuthHandlers(svc)
	r := gin.New()
	r.POST("/api/auth/login", h.Login)
	r.POST("/api/auth/logout", h.Logout)
	r.GET("/api/me", h.Me)
	return r
}

func TestLogin_Success(t *testing.T) {
	// Generate a valid bcrypt hash for the test password
	pwHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user := &entities.User{
		UserID:       1,
		Email:        "test@example.com",
		PasswordHash: string(pwHash),
		FirstName:    "Test",
		LastName:     "User",
	}
	userRepo := &fakeUserRepoAuth{user: user}
	cache := &fakeCacheAuth{}
	r := setupAuthRouter(userRepo, cache)

	body := `{"email":"test@example.com","password":"password"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	userRepo := &fakeUserRepoAuth{err: domain.ErrNotFound}
	cache := &fakeCacheAuth{}
	r := setupAuthRouter(userRepo, cache)

	body := `{"email":"test@example.com","password":"wrong"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	userRepo := &fakeUserRepoAuth{}
	cache := &fakeCacheAuth{}
	r := setupAuthRouter(userRepo, cache)

	body := `{"email":"test@example.com","password"`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestLogout_Success(t *testing.T) {
	userRepo := &fakeUserRepoAuth{}
	cache := &fakeCacheAuth{}
	r := setupAuthRouter(userRepo, cache)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestLogout_MissingHeader(t *testing.T) {
	userRepo := &fakeUserRepoAuth{}
	cache := &fakeCacheAuth{}
	r := setupAuthRouter(userRepo, cache)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestLogout_InvalidHeader(t *testing.T) {
	userRepo := &fakeUserRepoAuth{}
	cache := &fakeCacheAuth{}
	r := setupAuthRouter(userRepo, cache)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	req.Header.Set("Authorization", "Invalid token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestMe_WithAllFields(t *testing.T) {
	userRepo := &fakeUserRepoAuth{}
	cache := &fakeCacheAuth{}
	r := setupAuthRouter(userRepo, cache)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(123))
		c.Set("permissions", []string{"read", "write"})
		c.Set("first_name", "John")
		c.Set("last_name", "Doe")
		c.Set("email", "john@example.com")
		c.Next()
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestMe_WithMinimalFields(t *testing.T) {
	userRepo := &fakeUserRepoAuth{}
	cache := &fakeCacheAuth{}
	r := setupAuthRouter(userRepo, cache)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(123))
		c.Next()
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}
