package middlewares

import (
	"api/internal/application/useCase"
	"fmt"
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
		tokenStr, ok := bearerToken(c)
		if !ok {
			return
		}

		if abortIfInvalidated(c, authService, tokenStr) {
			return
		}

		claims := &useCase.AuthClaims{}
		if !parseTokenWithClaims(c, tokenStr, secret, claims) {
			return
		}

		if !validateTimeClaims(c, claims) {
			return
		}

		uid, ok := subjectToUserID(c, claims.Subject)
		if !ok {
			return
		}

		c.Set("userID", uid)
		c.Set("permissions", claims.Permissions)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// bearerToken extrae y valida el token Bearer del header Authorization.
func bearerToken(c *gin.Context) (string, bool) {
	auth := c.GetHeader("Authorization")
	parts := strings.Fields(auth)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		abortUnauthorized(c, "token requerido")
		return "", false
	}
	return parts[1], true
}

// abortIfInvalidated verifica si el token fue invalidado por el servicio.
func abortIfInvalidated(c *gin.Context, authService *useCase.AuthService, tokenStr string) bool {
	if authService.IsTokenInvalid(tokenStr) {
		abortUnauthorized(c, "token inválido")
		return true
	}
	return false
}

// parseTokenWithClaims parsea y valida el token JWT y su algoritmo.
func parseTokenWithClaims(c *gin.Context, tokenStr, secret string, claims *useCase.AuthClaims) bool {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("algoritmo inválido")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		abortUnauthorized(c, "token inválido")
		return false
	}
	return true
}

// validateTimeClaims valida expiración y not-before.
func validateTimeClaims(c *gin.Context, claims *useCase.AuthClaims) bool {
	now := time.Now()
	if claims.ExpiresAt != nil && now.After(claims.ExpiresAt.Time) {
		abortUnauthorized(c, "token expirado")
		return false
	}
	if claims.NotBefore != nil && now.Before(claims.NotBefore.Time) {
		abortUnauthorized(c, "token aún no válido")
		return false
	}
	return true
}

// subjectToUserID valida y convierte el Subject a userID.
func subjectToUserID(c *gin.Context, sub string) (uint, bool) {
	if sub == "" {
		abortUnauthorized(c, "sub faltante")
		return 0, false
	}
	uid64, err := strconv.ParseUint(sub, 10, 64)
	if err != nil || uid64 == 0 {
		abortUnauthorized(c, "sub inválido")
		return 0, false
	}
	return uint(uid64), true
}

// abortUnauthorized centraliza la respuesta 401 con mensaje de error.
func abortUnauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msg})
}
