package handlers

import (
	"errors"
	"net/http"

	"api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/requests"
	"github.com/gin-gonic/gin"
)

type UploadsHandlers struct {
	uploadsUC *useCase.UploadsUseCase
}

func NewUploadsHandlers(uploadsUC *useCase.UploadsUseCase) *UploadsHandlers {
	return &UploadsHandlers{uploadsUC: uploadsUC}
}

// CreatePostPolicy handles POST /api/uploads and returns a signed S3 POST policy for direct uploads.
func (h *UploadsHandlers) CreatePostPolicy(c *gin.Context) {
	var req requests.CreateUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	resp, err := h.uploadsUC.CreatePostPolicy(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}
