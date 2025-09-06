package handlers

import (
	"api/internal/application/useCase"
	"api/internal/domain/entities"
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type VideoHandlers struct {
	uploadsUC *useCase.UploadsUseCase
}

func NewVideoHandlers(uploadsUC *useCase.UploadsUseCase) *VideoHandlers {
	return &VideoHandlers{uploadsUC: uploadsUC}
}

func (h *VideoHandlers) Upload(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video file is required"})
		return
	}

	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID missing in context"})
		return
	}
	userID, ok := uidVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userID in context"})
		return
	}

	permsVal, ok := c.Get("permissions")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "permissions missing in context"})
		return
	}
	perms, ok := permsVal.([]string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid permissions in context"})
		return
	}
	allowed := slices.Contains(perms, "upload_video")
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	status := c.PostForm("status_id")
	// Default to UPLOADED when not provided; otherwise validate/normalize
	if status == "" {
		status = string(entities.StatusUploaded)
	} else {
		status = strings.ToUpper(status)
		valid := false
		for _, st := range entities.AllVideoStatuses() {
			if string(st) == status {
				valid = true
				break
			}
		}
		if !valid {
			// Return 400 with allowed statuses for clarity
			allowed := make([]string, 0)
			for _, st := range entities.AllVideoStatuses() {
				allowed = append(allowed, string(st))
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status", "allowed": allowed})
			return
		}
	}

	input := useCase.UploadVideoInput{
		Title:      title,
		FileHeader: file,
		Status:     status,
	}
	ctx := context.WithValue(c.Request.Context(), useCase.UserIDContextKey, userID)
	output, err := h.uploadsUC.UploadMultipart(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
