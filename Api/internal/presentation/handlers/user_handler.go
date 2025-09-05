package handlers

import (
	"errors"
	"main_videork/internal/application/useCase"
	"main_videork/internal/domain"
	"main_videork/internal/domain/requests"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	service *useCase.UserService
}

func NewUserHandlers(service *useCase.UserService) *UserHandlers {
	return &UserHandlers{service: service}
}

func (handler *UserHandlers) Register(context *gin.Context) {
	var request requests.RegisterUserRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := strings.ToLower(strings.TrimSpace(request.Email))

	if _, err := handler.service.GetByEmail(context.Request.Context(), email); err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "email_already_in_use"})
		return
	} else if !errors.Is(err, domain.ErrNotFound) {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if request.Password1 != request.Password2 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "passwords_do_not_match"})
		return
	}

	user, err := handler.service.CreateUser(
		context.Request.Context(),
		request.FirstName,
		request.LastName,
		email,
		request.Password1,
		request.CityID,
	)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"id":         user.UserID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"city_id":    user.CityID,
	})
}
