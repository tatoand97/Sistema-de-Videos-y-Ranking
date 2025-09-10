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

type VideoHandlers struct{ uploadsUC *useCase.UploadsUseCase }

func NewVideoHandlers(uploadsUC *useCase.UploadsUseCase) *VideoHandlers {
	return &VideoHandlers{uploadsUC: uploadsUC}
}

// ListVideos handles GET /api/videos (authenticated)
func (h *VideoHandlers) ListVideos(c *gin.Context) {
	userID, ok := userIDFromContextOrAbort(c)
	if !ok {
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
		out = append(out, toVideoResponse(v))
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

	// Validar tamaño declarado (=100MB)
	if file.Size > validations.MaxBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 100MB)"})
		return
	}

	userID, ok := userIDFromContextOrAbort(c)
	if !ok {
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
	// Align with OpenAPI: expect "videos:upload".
	// Keep backward compatibility with legacy "upload_video".
	allowed := slices.Contains(perms, "videos:upload") || slices.Contains(perms, "upload_video")
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
	userID, ok := userIDFromContextOrAbort(c)
	if !ok {
		return
	}

	// Path param
	vidUint, ok := parseVideoIDOrAbort(c)
	if !ok {
		return
	}

	// Query use case enforcing ownership
	v, err := h.uploadsUC.GetUserVideoByID(c.Request.Context(), userID, vidUint)
	if err != nil {
		if handled := writeStandardDomainError(c, err); handled {
			return
		}
	}

	resp := toVideoResponse(v)
	c.JSON(http.StatusOK, resp)
}

// DeleteVideo handles DELETE /api/videos/:video_id (authenticated, own video only)
func (h *VideoHandlers) DeleteVideo(c *gin.Context) {
	// 1) Auth
	userID, ok := userIDFromContextOrAbort(c)
	if !ok {
		return
	}

	// 2) Path param
	vidStr := c.Param("video_id")
	vidUint, ok := parseVideoIDOrAbort(c)
	if !ok {
		return
	}

	// 3-6) Execute delete with eligibility rules in use case
	err := h.uploadsUC.DeleteUserVideoIfEligible(c.Request.Context(), userID, vidUint)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalid):
			// Elegibilidad incumplida (p.ej., publicado para votación o procesado)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "No se puede eliminar: el video está publicado para votación o ya fue procesado."})
			return
		default:
			if handled := writeStandardDomainError(c, err); handled {
				return
			}
		}
	}

	// 7) Success 200 with exact body
	c.JSON(http.StatusOK, gin.H{
		"message":  "El video ha sido eliminado exitosamente.",
		"video_id": vidStr,
	})
}

// PublishVideo handles POST /api/videos/:video_id/publish (moderation action)
// Requires permission: edit_video or moderate_content
func (h *VideoHandlers) PublishVideo(c *gin.Context) {
	// 1) Auth
	if _, ok := userIDFromContextOrAbort(c); !ok {
		return
	}
	// 2) Permissions
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
	allowed := slices.Contains(perms, "edit_video") || slices.Contains(perms, "moderate_content")
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// 3) Path param
	vidUint, ok := parseVideoIDOrAbort(c)
	if !ok {
		return
	}

	// 4) Publish via use case
	if err := h.uploadsUC.PublishVideo(c.Request.Context(), vidUint); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Video no encontrado."})
			return
		case errors.Is(err, domain.ErrInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "El video no está listo para publicarse."})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video publicado exitosamente.", "video_id": c.Param("video_id")})
}

// --- Helpers to reduce duplication ---

// userIDFromContextOrAbort extracts userID from context or writes a 401 and returns false.
func userIDFromContextOrAbort(c *gin.Context) (uint, bool) {
	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": invalidTokenExpiredMsg})
		return 0, false
	}
	userID, ok := uidVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": invalidTokenExpiredMsg})
		return 0, false
	}
	return userID, true
}

// parseVideoIDOrAbort validates path param "video_id" and returns it as uint.
func parseVideoIDOrAbort(c *gin.Context) (uint, bool) {
	vidStr := c.Param("video_id")
	if vidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Parámetro inválido."})
		return 0, false
	}
	parsed, err := strconv.ParseUint(vidStr, 10, 64)
	if err != nil || parsed == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Parámetro inválido."})
		return 0, false
	}
	return uint(parsed), true
}

// toVideoResponse maps an entities.Video to responses.VideoResponse.
func toVideoResponse(v *entities.Video) responses.VideoResponse {
	status := "uploaded"
	if v.Status == string(entities.StatusPublished) {
		status = "published"
	} else if v.Status == string(entities.StatusProcessed) {
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
	return resp
}

// writeStandardDomainError writes common domain error translations.
// Returns true if it wrote a response.
func writeStandardDomainError(c *gin.Context, err error) bool {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Video no encontrado."})
		return true
	case errors.Is(err, domain.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "Acceso denegado."})
		return true
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return true
	}
}
