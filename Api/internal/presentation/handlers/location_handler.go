package handlers

import (
	"main_videork/internal/application/useCase"
	"main_videork/internal/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type LocationHandlers struct {
	service *useCase.LocationService
}

func NewLocationHandlers(service *useCase.LocationService) *LocationHandlers {
	return &LocationHandlers{service: service}
}

func (h *LocationHandlers) GetCityID(c *gin.Context) {
	country := strings.TrimSpace(c.Query("country"))
	city := strings.TrimSpace(c.Query("city"))
	if country == "" || city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "country and city are required"})
		return
	}

	id, err := h.service.GetCityID(c.Request.Context(), country, city)
	if err != nil {
		switch err {
		case domain.ErrInvalid:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid country or city"})
		case domain.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "city not found for country"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"city_id": id})
}
