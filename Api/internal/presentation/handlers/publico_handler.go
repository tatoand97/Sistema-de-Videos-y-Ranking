package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main_videork/internal/application/useCase"
)

// PublicHandlers maneja endpoints públicos relacionados a videos.
type PublicHandlers struct {
	service *useCase.PublicService
}

func NewPublicHandlers(service *useCase.PublicService) *PublicHandlers {
	return &PublicHandlers{service: service}
}

// ListPublicVideos maneja GET /api/public/videos
func (h *PublicHandlers) ListPublicVideos(c *gin.Context) {
	results, err := h.service.ListPublicVideos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}
