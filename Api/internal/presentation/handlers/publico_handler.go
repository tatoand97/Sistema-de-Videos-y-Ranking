package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"api/internal/application/useCase"
	"api/internal/domain"
	domainresponses "api/internal/domain/responses"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "No se pudo registrar el voto."})
		return
	}

	// 6) OK
	c.JSON(http.StatusOK, gin.H{"message": "Voto registrado exitosamente."})
}

// Unique violation mapping moved to repository layer; handler sees domain.ErrConflict only.

// ListRankings maneja GET /api/public/rankings
// Público, sin autenticación. Devuelve un array de RankingEntry.
func (h *PublicHandlers) ListRankings(c *gin.Context) {
	// Validación de parámetros de consulta
	cityParam := strings.TrimSpace(c.Query("city"))
	var city *string
	if cityParam != "" {
		city = &cityParam
	}

	// Defaults
	page := 1
	pageSize := 20

	if ps := strings.TrimSpace(c.Query("page")); ps != "" {
		v, err := strconv.Atoi(ps)
		if err != nil || v < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Parámetro inválido en la consulta"})
			return
		}
		page = v
	}
	if pss := strings.TrimSpace(c.Query("pageSize")); pss != "" {
		v, err := strconv.Atoi(pss)
		if err != nil || v < 1 || v > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Parámetro inválido en la consulta"})
			return
		}
		pageSize = v
	}

	items, err := h.service.Rankings(c.Request.Context(), city, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}

	// Mapear a salida pública con posición por página (empezando en 1)
	resp := make([]domainresponses.RankingEntry, 0, len(items))
	for i, it := range items {
		resp = append(resp, domainresponses.RankingEntry{
			Position: i + 1,
			Username: it.Username,
			City:     it.City,
			Votes:    it.Votes,
		})
	}
	// Sin headers adicionales de paginación no normados.
	c.JSON(http.StatusOK, resp)
}
