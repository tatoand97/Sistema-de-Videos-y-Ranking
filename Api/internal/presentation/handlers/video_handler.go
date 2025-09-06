package handlers

import (
	"api/internal/application/useCase"
	"api/internal/application/validations"
	"api/internal/domain"
	"api/internal/domain/entities"
	"api/internal/domain/responses"
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VideoHandlers struct {
	uploadsUC *useCase.UploadsUseCase
}

func NewVideoHandlers(uploadsUC *useCase.UploadsUseCase) *VideoHandlers {
	return &VideoHandlers{uploadsUC: uploadsUC}
}

// ListVideos handles GET /api/videos (authenticated)
func (h *VideoHandlers) ListVideos(c *gin.Context) {
	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Token inválido o expirado."})
		return
	}
	userID, ok := uidVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Token inválido o expirado."})
		return
	}

	videos, err := h.uploadsUC.ListUserVideos(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}

	// Map to OAS Video schema
	out := make([]responses.VideoResponse, 0, len(videos))
	for _, v := range videos {
		// Map status to storage status: uploaded | processed
		status := "uploaded"
		if v.Status == string(entities.StatusProcessed) {
			status = "processed"
		}

		var originalURL *string
		if v.OriginalFile != "" {
			s := v.OriginalFile
			originalURL = &s
		}

		vr := responses.VideoResponse{
			VideoID:     fmt.Sprintf("%d", v.VideoID),
			Title:       v.Title,
			Status:      status,
			UploadedAt:  v.UploadedAt,
			ProcessedAt: v.ProcessedAt,
			OriginalURL: originalURL,
		}
		if v.ProcessedFile != nil && *v.ProcessedFile != "" {
			vr.ProcessedURL = v.ProcessedFile
		}
		out = append(out, vr)
	}

	c.JSON(http.StatusOK, out)
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := uidVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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

// GetVideoDetail handles GET /api/videos/:video_id (authenticated, own video only)
func (h *VideoHandlers) GetVideoDetail(c *gin.Context) {
	// Auth identity from context
	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Token inválido o expirado."})
		return
	}
	userID, ok := uidVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Token inválido o expirado."})
		return
	}

	// Path param
	vidStr := c.Param("video_id")
	if vidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Parámetro inválido."})
		return
	}
	// video_id is numeric in our DB
	var vidUint uint
	{
		var parsed uint64
		var err error
		parsed, err = strconv.ParseUint(vidStr, 10, 64)
		if err != nil || parsed == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Parámetro inválido."})
			return
		}
		vidUint = uint(parsed)
	}

	// Query use case enforcing ownership
	v, err := h.uploadsUC.GetUserVideoByID(c.Request.Context(), userID, vidUint)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Video no encontrado."})
			return
		case errors.Is(err, domain.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "Acceso denegado."})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
			return
		}
	}

	// Map status: uploaded | processed
	status := "uploaded"
	if v.Status == string(entities.StatusProcessed) {
		status = "processed"
	}

	var originalURL *string
	if v.OriginalFile != "" {
		s := v.OriginalFile
		originalURL = &s
	}

	resp := responses.VideoResponse{
		VideoID:     fmt.Sprintf("%d", v.VideoID),
		Title:       v.Title,
		Status:      status,
		UploadedAt:  v.UploadedAt,
		ProcessedAt: v.ProcessedAt,
		OriginalURL: originalURL,
	}
	if v.ProcessedFile != nil && *v.ProcessedFile != "" {
		resp.ProcessedURL = v.ProcessedFile
	}
	c.JSON(http.StatusOK, resp)
}
