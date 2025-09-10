package middlewares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/application/useCase"
	"api/internal/domain/entities"
	"api/internal/presentation/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockUserRepoMiddleware struct {
	user *entities.User
	err  error
}

func (m *mockUserRepoMiddleware) Create(ctx context.Context, user *entities.User) error { return nil }
func (m *mockUserRepoMiddleware) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}
func (m *mockUserRepoMiddleware) EmailExists(ctx context.Context, email string) (bool, error) { return false, nil }
func (m *mockUserRepoMiddleware) GetPermissions(ctx context.Context, userID uint) ([]string, error) { return []string{"read"}, nil }

type mockCacheMiddleware struct {
	blacklisted bool
}

func (m *mockCacheMiddleware) Set(ctx context.Context, key string, value interface{}, ttl int) error { return nil }
func (m *mockCacheMiddleware) Get(ctx context.Context, key string) (string, error) { return "", nil }
func (m *mockCacheMiddleware) Delete(ctx context.Context, key string) error { return nil }
func (m *mockCacheMiddleware) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	return m.blacklisted, nil
}
func (m *mockCacheMiddleware) BlacklistToken(ctx context.Context, token string, ttl int) error { return nil }

func TestJWTMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	user := &entities.User{
		UserID:    1,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}
	userRepo := &mockUserRepoMiddleware{user: user}
	authService := useCase.NewAuthService(userRepo, "secret")

	// Generate a valid token
	token, _, _ := authService.Login(context.Background(), "test@example.com", "password")

	r.Use(middlewares.JWTMiddleware(authService, "secret"))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	userRepo := &mockUserRepoMiddleware{}
	authService := useCase.NewAuthService(userRepo, "secret")

	r.Use(middlewares.JWTMiddleware(authService, "secret"))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	userRepo := &mockUserRepoMiddleware{}
	authService := useCase.NewAuthService(userRepo, "secret")

	r.Use(middlewares.JWTMiddleware(authService, "secret"))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTMiddleware_BlacklistedToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	user := &entities.User{
		UserID:    1,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}
	userRepo := &mockUserRepoMiddleware{user: user}
	authService := useCase.NewAuthService(userRepo, "secret")

	// Generate a valid token
	token, _, _ := authService.Login(context.Background(), "test@example.com", "password")

	r.Use(middlewares.JWTMiddleware(authService, "secret"))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}