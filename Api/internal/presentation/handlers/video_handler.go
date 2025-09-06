package handlers

import (
	"api/internal/application/useCase"
	"api/internal/application/validations"
	"api/internal/domain"
	"api/internal/domain/entities"
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// Acepta el nombre estándar del documento: video_file. Mantiene compatibilidad con "video".
	file, err := c.FormFile("video_file")
	if err != nil {
		file, err = c.FormFile("video")
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video file is required"})
		return
	}

	// Validar MIME declarado
	if ct := file.Header.Get("Content-Type"); ct != "video/mp4" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mimeType must be video/mp4"})
		return
	}

	// Validar tamaño declarado (≤100MB)
	if file.Size > validations.MaxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 100MB)"})
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

	// Forzar estado UPLOADED según documento (no se controla por API)
	status := string(entities.StatusUploaded)

	input := useCase.UploadVideoInput{
		Title:      title,
		FileHeader: file,
		Status:     status,
	}
	ctx := context.WithValue(c.Request.Context(), useCase.UserIDContextKey, userID)
	output, err := h.uploadsUC.UploadMultipart(ctx, input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_ = output // No se expone el Video en 201
	c.JSON(http.StatusCreated, gin.H{
		"message": "Video subido correctamente. Procesamiento en curso.",
		"task_id": uuid.NewString(),
	})
}
