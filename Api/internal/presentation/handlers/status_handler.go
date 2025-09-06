package handlers

import (
	"api/internal/application/useCase"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusHandlers struct {
	statusService *useCase.StatusService
}

func NewStatusHandlers(statusService *useCase.StatusService) *StatusHandlers {
	return &StatusHandlers{statusService: statusService}
}

// ListVideoStatuses handles GET /api/videos/statuses
func (h *StatusHandlers) ListVideoStatuses(c *gin.Context) {
	_ = context.Background() // kept for future parity; not necessary now
	statuses := h.statusService.ListVideoStatuses(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"statuses": statuses})
}
