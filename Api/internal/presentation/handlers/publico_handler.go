package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"main_videork/internal/application/useCase"
	"main_videork/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

// PublicHandlers maneja endpoints públicos relacionados a videos.
type PublicHandlers struct {
	service *useCase.PublicService
}

func NewPublicHandlers(service *useCase.PublicService) *PublicHandlers {
	return &PublicHandlers{service: service}
}

// ListPublicVideos maneja GET /api/public/videos
func (h *PublicHandlers) ListPublicVideos(c *gin.Context) {
	results, err := h.service.ListPublicVideos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// VotePublicVideo maneja POST /api/public/videos/:video_id/vote
func (h *PublicHandlers) VotePublicVideo(c *gin.Context) {
	// 1) Auth: extraer userID del contexto
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Token inválido o expirado."})
		return
	}

	// 2) Path param
	v := c.Param("video_id")
	if v == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Parámetros inválidos."})
		return
	}
	vid64, err := strconv.ParseUint(v, 10, 64)
	if err != nil || vid64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Parámetros inválidos."})
		return
	}
	videoID := uint(vid64)

	// 3-5) Lógica de voto a través del servicio (incluye verificación de existencia y unicidad)
	err = h.service.VotePublicVideo(c.Request.Context(), videoID, userID)
	if err != nil {
		// Mapear errores comunes
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Video no encontrado."})
			return
		}
		if errors.Is(err, domain.ErrConflict) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Ya has votado por este video."})
			return
		}
		if isUniqueViolation(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Ya has votado por este video."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "No se pudo registrar el voto."})
		return
	}

	// 6) OK
	c.JSON(http.StatusOK, gin.H{"message": "Voto registrado exitosamente."})
}

// isUniqueViolation detecta error 23505 de Postgres
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return true
		}
	}
	return false
}
