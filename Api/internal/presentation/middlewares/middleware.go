package middlewares

import (
	"fmt"
	"main_videork/internal/application/useCase"
	"main_videork/internal/domain/interfaces"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware valida el token Bearer, verifica expiración y not-before,
// y coloca el userID en el contexto para el resto de handlers protegidos.
func JWTMiddleware(authService *useCase.AuthService, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		parts := strings.Fields(auth)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token requerido"})
			return
		}
		tokenStr := parts[1]

		if authService.IsTokenInvalid(tokenStr) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("algoritmo inválido")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}

		now := time.Now()
		// Verifica expiración
		if claims.ExpiresAt != nil && now.After(claims.ExpiresAt.Time) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expirado"})
			return
		}
		// Verifica not-before
		if claims.NotBefore != nil && now.Before(claims.NotBefore.Time) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token aún no válido"})
			return
		}

		// Extraer userID desde Subject
		if claims.Subject == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sub faltante"})
			return
		}
		uid64, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil || uid64 == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sub inválido"})
			return
		}
		c.Set("userID", uint(uid64))

		c.Next()
	}
}

// PermissionMiddleware verifies that the authenticated user has all required permissions.
func PermissionMiddleware(userRepo interfaces.UserRepository, required ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uidVal, ok := c.Get("userID")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "userID missing"})
			return
		}
		uid, ok := uidVal.(uint)
		if !ok || uid == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "userID invalid"})
			return
		}

		perms, err := userRepo.GetPermissions(c.Request.Context(), uid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot load permissions"})
			return
		}
		permSet := make(map[string]struct{}, len(perms))
		for _, p := range perms {
			permSet[p] = struct{}{}
		}
		for _, rp := range required {
			if _, ok := permSet[rp]; !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
		}
		c.Next()
	}
}
