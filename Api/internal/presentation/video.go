package presentation

import (
	"main_videork/internal/application/useCase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// VideoHandlers manages video related endpoints.
type VideoHandlers struct {
	uploadUC *useCase.UploadVideoUseCase
}

// NewVideoHandlers creates a new VideoHandlers instance.
func NewVideoHandlers(uploadUC *useCase.UploadVideoUseCase) *VideoHandlers {
	return &VideoHandlers{uploadUC: uploadUC}
}

// Upload handles receiving a video file and a title via multipart form.
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

	playerIDStr := c.PostForm("player_id")
	playerID, err := strconv.ParseUint(playerIDStr, 10, 64)
	if err != nil || playerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "player_id is required and must be a valid uint"})
		return
	}

	statusIDStr := c.PostForm("status_id")
	statusID, err := strconv.ParseUint(statusIDStr, 10, 64)
	if err != nil || statusID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status_id is required and must be a valid uint"})
		return
	}

	input := useCase.UploadVideoInput{
		PlayerID:   uint(playerID),
		Title:      title,
		FileHeader: file,
		StatusID:   uint(statusID),
	}
	output, err := h.uploadUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
