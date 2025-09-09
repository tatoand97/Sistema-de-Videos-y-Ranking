package handlers

import (
	"api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/requests"
	"errors"
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
		context.JSON(http.StatusBadRequest, gin.H{"error": badRequest, "message": err.Error()})
		return
	}

	email := strings.ToLower(strings.TrimSpace(request.Email))

	if exists, err := handler.service.EmailExists(context.Request.Context(), email); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	} else if exists {
		context.JSON(http.StatusBadRequest, gin.H{"error": badRequest, "message": "email already in use"})
		return
	}

	if request.Password1 != request.Password2 {
		context.JSON(http.StatusBadRequest, gin.H{"error": badRequest, "message": "passwords do not match"})
		return
	}

	user, err := handler.service.CreateUser(
		context.Request.Context(),
		request.FirstName,
		request.LastName,
		email,
		request.Password1,
		request.Country,
		request.City,
	)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) || errors.Is(err, domain.ErrInvalid) {
			context.JSON(http.StatusBadRequest, gin.H{"error": badRequest, "message": "invalid city or country"})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}

	_ = user // no se expone en la respuesta según OpenAPI
	context.JSON(http.StatusCreated, gin.H{
		"message": "Usuario creado exitosamente.",
	})
}
