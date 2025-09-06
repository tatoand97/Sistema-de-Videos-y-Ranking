package handlers

import (
	"main_videork/internal/application/useCase"
	"main_videork/internal/domain/requests"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandlers struct {
	service *useCase.AuthService
}

func NewAuthHandlers(service *useCase.AuthService) *AuthHandlers {
	return &AuthHandlers{service: service}
}

func (handler *AuthHandlers) Login(context *gin.Context) {
	var request requests.LoginRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, expiresIn, err := handler.service.Login(context.Request.Context(), request.Email, request.Password)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"token_type":   "Bearer",
		"expires_in":   expiresIn,
		"access_token": token,
	})
}

func (handler *AuthHandlers) Logout(context *gin.Context) {
	header := context.GetHeader("Authorization")
	const prefix = "Bearer "
	if header == "" || !strings.HasPrefix(header, prefix) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		return
	}
	token := strings.TrimSpace(header[len(prefix):])
	if token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	if err := handler.service.Logout(token); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.Status(http.StatusNoContent)
}

func (handler *AuthHandlers) Me(context *gin.Context) {
	uidVal, _ := context.Get("userID")
	permsVal, _ := context.Get("permissions")
	firstNameVal, _ := context.Get("first_name")
	lastNameVal, _ := context.Get("last_name")
	emailVal, _ := context.Get("email")

	response := gin.H{
		"status": "ok",
	}

	if uid, ok := uidVal.(uint); ok {
		response["user_id"] = uid
	}
	if perms, ok := permsVal.([]string); ok {
		response["permissions"] = perms
	}
	if fn, ok := firstNameVal.(string); ok && fn != "" {
		response["first_name"] = fn
	}
	if ln, ok := lastNameVal.(string); ok && ln != "" {
		response["last_name"] = ln
	}
	if em, ok := emailVal.(string); ok && em != "" {
		response["email"] = em
	}

	context.JSON(http.StatusOK, response)
}
