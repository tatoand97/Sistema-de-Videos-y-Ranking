package presentation

import (
	"main_videork/internal/application/useCase"
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

type registerRequest struct {
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`
	Email     string `json:"email" form:"email" binding:"required"`
	Password1 string `json:"password1" form:"password1" binding:"required"`
	Password2 string `json:"password2" form:"password2" binding:"required"`
	City      string `json:"city" form:"city"`
	Country   string `json:"country" form:"country"`
}

func (handler *AuthHandlers) Register(context *gin.Context) {
	var request registerRequest
	if err := context.ShouldBind(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := strings.ToLower(strings.TrimSpace(request.Email))

	exists, err := handler.service.EmailExists(context.Request.Context(), email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		context.JSON(http.StatusConflict, gin.H{"error": "email_already_in_use"})
		return
	}

	if request.Password1 != request.Password2 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "passwords_do_not_match"})
		return
	}

	user, err := handler.service.Register(context.Request.Context(), request.FirstName, request.LastName, request.Email, request.Password1)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"id":         user.UserId,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"password":   user.PasswordHash,
	})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (handler *AuthHandlers) Login(context *gin.Context) {
	var request loginRequest
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
	if err := handler.service.Logout(context.Request.Context(), token); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.Status(http.StatusNoContent)
}
