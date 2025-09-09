package handlers_test

import (
	"api/internal/application/useCase"
	"api/internal/presentation/handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandlers_Login_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Minimal service; won't be reached due to invalid JSON
	svc := useCase.NewAuthService(nil, "secret")
	h := handlers.NewAuthHandlers(svc)

	r.POST("/api/auth/login", h.Login)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandlers_Logout_InvalidHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := useCase.NewAuthService(nil, "secret")
	h := handlers.NewAuthHandlers(svc)
	r.POST("/api/auth/logout", h.Logout)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandlers_Me_PropagatesContextValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := useCase.NewAuthService(nil, "secret")
	h := handlers.NewAuthHandlers(svc)
	r.GET("/api/me", func(c *gin.Context) {
		c.Set("userID", uint(42))
		c.Set("permissions", []string{"read"})
		c.Set("first_name", "Ada")
		c.Set("last_name", "Lovelace")
		c.Set("email", "ada@example.com")
		h.Me(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"status\":\"ok\"")
	assert.Contains(t, w.Body.String(), "\"user_id\":42")
	assert.Contains(t, w.Body.String(), "\"permissions\":[\"read\"]")
	assert.Contains(t, w.Body.String(), "\"first_name\":\"Ada\"")
	assert.Contains(t, w.Body.String(), "\"email\":\"ada@example.com\"")
}
