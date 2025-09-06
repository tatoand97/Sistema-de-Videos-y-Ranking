package handlers

import (
	"context"
	"main_videork/internal/application/useCase"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type VideoHandlers struct {
	uploadUC *useCase.UploadVideoUseCase
}

func NewVideoHandlers(uploadUC *useCase.UploadVideoUseCase) *VideoHandlers {
	return &VideoHandlers{uploadUC: uploadUC}
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

	input := useCase.UploadVideoInput{
		Title:      title,
		FileHeader: file,
		Status:     status,
	}
	ctx := context.WithValue(c.Request.Context(), useCase.UserIDContextKey, userID)
	output, err := h.uploadUC.Execute(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
