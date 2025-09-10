package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"api/internal/application/useCase"
	"api/internal/domain"
	"github.com/gin-gonic/gin"
)

type UploadsHandlers struct {
	uploadsUC *useCase.UploadsUseCase
}

func NewUploadsHandlers(uploadsUC *useCase.UploadsUseCase) *UploadsHandlers {
	return &UploadsHandlers{uploadsUC: uploadsUC}
}

// NewUploadsHandler provides backward-compatible constructor name expected by some tests.
func NewUploadsHandler(uploadsUC *useCase.UploadsUseCase) *UploadsHandlers {
	return NewUploadsHandlers(uploadsUC)
}

// UploadVideo handles POST /api/uploads multipart uploads.
// It bridges Gin context userID into the use case context and delegates to UploadMultipart.
func (h *UploadsHandlers) UploadVideo(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	status := c.PostForm("status")
	if status == "" {
		status = "UPLOADED"
	}

	// Accept standard form field name "file"; keep compatibility with "video" and "video_file".
	fh, err := c.FormFile("file")
	if err != nil {
		fh, err = c.FormFile("video")
	}
	if err != nil {
		fh, err = c.FormFile("video_file")
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video file is required"})
		return
	}

	// Extract userID from Gin context
	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := uidVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Inject userID into request context as expected by use case
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, useCase.UserIDContextKey, userID)

	out, err := h.uploadsUC.UploadMultipart(ctx, useCase.UploadVideoInput{
		Title:      title,
		FileHeader: fh,
		Status:     status,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"video_id": fmt.Sprintf("%d", out.VideoID),
		"title":    out.Title,
	})
}
