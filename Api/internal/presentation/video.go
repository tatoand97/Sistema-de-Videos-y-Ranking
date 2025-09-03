package presentation

import (
	"context"
	"main_videork/internal/application/useCase"
	"net/http"
	"strconv"

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

	statusIDStr := c.PostForm("status_id")
	statusID, err := strconv.ParseUint(statusIDStr, 10, 64)
	if err != nil || statusID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status_id is required and must be a valid uint"})
		return
	}

	input := useCase.UploadVideoInput{
		Title:      title,
		FileHeader: file,
		StatusID:   uint(statusID),
	}
	ctx := context.WithValue(c.Request.Context(), useCase.UserIDContextKey, userID)
	output, err := h.uploadUC.Execute(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
