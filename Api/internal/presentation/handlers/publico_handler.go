package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/interfaces"
	domainresponses "api/internal/domain/responses"

	"github.com/gin-gonic/gin"
)

// PublicHandlers maneja endpoints publicos relacionados a videos.
type PublicHandlers struct {
	service        *useCase.PublicService
	cache          interfaces.Cache
	idemTTLSeconds int
}

// NewPublicHandlers mantiene compatibilidad para tests y uso sin cache/idempotencia.
func NewPublicHandlers(service *useCase.PublicService) *PublicHandlers {
	return &PublicHandlers{service: service}
}

// NewPublicHandlersWithCache permite inyectar cache y TTL de idempotencia.
func NewPublicHandlersWithCache(service *useCase.PublicService, cache interfaces.Cache, idemTTLSeconds int) *PublicHandlers {
	return &PublicHandlers{service: service, cache: cache, idemTTLSeconds: idemTTLSeconds}
}

// ListPublicVideos maneja GET /api/public/videos
func (h *PublicHandlers) ListPublicVideos(c *gin.Context) {
	results, err := h.service.ListPublicVideos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// VotePublicVideo maneja POST /api/public/videos/:video_id/vote
func (h *PublicHandlers) VotePublicVideo(c *gin.Context) {
	// 1) Auth: extraer userID del contexto
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": invalidTokenExpiredMsg})
		return
	}

	// 2) Path param
	v := c.Param("video_id")
	if v == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": badRequest, "message": "Parametros invalidos."})
		return
	}
	vid64, err := strconv.ParseUint(v, 10, 64)
	if err != nil || vid64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": badRequest, "message": "Parametros invalidos."})
		return
	}
	videoID := uint(vid64)

	// 3) Idempotencia opcional: X-Event-Id (header) o query param "eventId"
	var eventIDPtr *string
	if evt := strings.TrimSpace(c.GetHeader("X-Event-Id")); evt != "" {
		eventIDPtr = &[]string{evt}[0]
	} else if evtq := strings.TrimSpace(c.Query("eventId")); evtq != "" {
		eventIDPtr = &[]string{evtq}[0]
	}
	// Idempotencia rapida con Redis: degradar si falla y liberar clave si BD falla
	var seenKey string
	nxSet := false
	if eventIDPtr != nil && h.cache != nil && h.idemTTLSeconds > 0 {
		seenKey = fmt.Sprintf("seen:{events:%s}", *eventIDPtr)
		ok, derr := h.cache.SetNX(c.Request.Context(), seenKey, []byte("1"), time.Duration(h.idemTTLSeconds)*time.Second)
		if derr == nil && !ok {
			// Evento ya procesado recientemente: exito idempotente
			c.JSON(http.StatusOK, gin.H{"message": "Voto registrado exitosamente."})
			return
		}
		if derr == nil && ok {
			nxSet = true
		}
		// Si derr != nil, seguimos y confiamos en la BD
	}

	// 4-6) Logica de voto via servicio (incluye verificacion de existencia y unicidad)
	if eventIDPtr != nil {
		err = h.service.VotePublicVideoWithEvent(c.Request.Context(), videoID, userID, eventIDPtr)
	} else {
		err = h.service.VotePublicVideo(c.Request.Context(), videoID, userID)
	}
	if err != nil {
		// Si fallo BD y marcamos NX, liberar para permitir reintentos
		if nxSet && h.cache != nil && seenKey != "" {
			_ = h.cache.DeleteWildcard(c.Request.Context(), seenKey)
		}
		// Mapear errores comunes
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Video no encontrado."})
			return
		}
		if errors.Is(err, domain.ErrIdempotent) {
			c.JSON(http.StatusOK, gin.H{"message": "Voto registrado exitosamente."})
			return
		}
		if errors.Is(err, domain.ErrConflict) {
			c.JSON(http.StatusBadRequest, gin.H{"error": badRequest, "message": "Ya has votado por este video."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "No se pudo registrar el voto."})
		return
	}

	// 5) Agregados Redis removidos; BD es fuente de verdad.

	// 6) OK
	c.JSON(http.StatusOK, gin.H{"message": "Voto registrado exitosamente."})
}

// ListRankings maneja GET /api/public/rankings
// Publico, sin autenticacion. Devuelve un array de RankingEntry.
func (h *PublicHandlers) ListRankings(c *gin.Context) {
	// Validacion de parametros de consulta
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
			c.JSON(http.StatusBadRequest, gin.H{"error": badRequest, "message": "Parametro invalido en la consulta"})
			return
		}
		page = v
	}
	if pss := strings.TrimSpace(c.Query("pageSize")); pss != "" {
		v, err := strconv.Atoi(pss)
		if err != nil || v < 1 || v > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": badRequest, "message": "Parametro invalido en la consulta"})
			return
		}
		pageSize = v
	}

	// Orquestación en el caso de uso (Redis -> fallback BD)
	items, err := h.service.Rankings(c.Request.Context(), city, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	resp := make([]domainresponses.RankingEntry, 0, len(items))
	for i, it := range items {
		resp = append(resp, domainresponses.RankingEntry{
			Position: i + 1,
			Username: it.Username,
			City:     it.City,
			Votes:    it.Votes,
		})
	}
	c.JSON(http.StatusOK, resp)
}
