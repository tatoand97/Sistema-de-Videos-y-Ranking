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
	agg            interfaces.Aggregates
}

// NewPublicHandlers mantiene compatibilidad para tests y uso sin cache/idempotencia.
func NewPublicHandlers(service *useCase.PublicService) *PublicHandlers {
	return &PublicHandlers{service: service}
}

// NewPublicHandlersWithCache permite inyectar cache y TTL de idempotencia.
func NewPublicHandlersWithCache(service *useCase.PublicService, cache interfaces.Cache, idemTTLSeconds int) *PublicHandlers {
	return &PublicHandlers{service: service, cache: cache, idemTTLSeconds: idemTTLSeconds}
}

// NewPublicHandlersFull allows injecting cache (for idempotency) and aggregates (for leaderboards/stats).
func NewPublicHandlersFull(service *useCase.PublicService, cache interfaces.Cache, idemTTLSeconds int, agg interfaces.Aggregates) *PublicHandlers {
	return &PublicHandlers{service: service, cache: cache, idemTTLSeconds: idemTTLSeconds, agg: agg}
}

// normalizeCityKey crea una clave segura (slug ASCII) para usar en Redis a partir del nombre de ciudad.
// - minúsculas
// - sin tildes/ñ
// - espacios a '-'
// - solo [a-z0-9_-]
func normalizeCityKey(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return s
	}
	r := strings.NewReplacer(
		"á", "a", "à", "a", "ä", "a", "â", "a", "ã", "a",
		"é", "e", "è", "e", "ë", "e", "ê", "e",
		"í", "i", "ì", "i", "ï", "i", "î", "i",
		"ó", "o", "ò", "o", "ö", "o", "ô", "o", "õ", "o",
		"ú", "u", "ù", "u", "ü", "u", "û", "u",
		"ñ", "n",
	)
	s = r.Replace(s)
	// Reemplazar espacios por '-'
	s = strings.ReplaceAll(s, " ", "-")
	// Filtrar caracteres no deseados
	b := make([]rune, 0, len(s))
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' {
			b = append(b, ch)
		}
	}
	return string(b)
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

	// 5) Actualizar agregados en Redis (best-effort, la BD es fuente de verdad)
	if h.agg != nil {
		nowUnix := time.Now().Unix()
		uniq := strconv.FormatUint(uint64(userID), 10)
		// Ranking de videos (global y por ciudad)
		videoMember := strconv.FormatUint(uint64(videoID), 10)
		_ = h.agg.UpdateAfterVote(c.Request.Context(), "videos", videoMember, &uniq, nowUnix, true)
		var cityKey string
		var ownerUserID uint
		if vid, err := h.service.GetPublicByID(c.Request.Context(), videoID); err == nil && vid != nil {
			if vid.City != nil {
				cityKey = normalizeCityKey(*vid.City)
				if cityKey != "" {
					_ = h.agg.UpdateAfterVote(c.Request.Context(), "videos:city:"+cityKey, videoMember, &uniq, nowUnix, true)
				}
			}
			ownerUserID = vid.OwnerUserID
		}
		// Ranking de usuarios (global y por ciudad)
		if ownerUserID != 0 {
			userMember := strconv.FormatUint(uint64(ownerUserID), 10)
			_ = h.agg.UpdateAfterVote(c.Request.Context(), "users", userMember, &uniq, nowUnix, true)
			if cityKey != "" {
				_ = h.agg.UpdateAfterVote(c.Request.Context(), "users:city:"+cityKey, userMember, &uniq, nowUnix, true)
			}
		}
	}

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

// GetLeaderboard maneja GET /api/public/leaderboard/:poll_id
// Devuelve miembros y puntajes desde Redis (ZSET) paginados por rango start..stop.
func (h *PublicHandlers) GetLeaderboard(c *gin.Context) {
	if h.agg == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service Unavailable", "message": "Aggregates not configured"})
		return
	}
	pollId := strings.TrimSpace(c.Param("poll_id"))
	if pollId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": badRequest, "message": "poll_id requerido"})
		return
	}
	// Filtro opcional por ciudad: si viene, usamos <pollId>:city:<cityKey>
	if city := strings.TrimSpace(c.Query("city")); city != "" {
		cityKey := normalizeCityKey(city)
		// Evitar duplicar sufijo si ya vino formado
		if !strings.Contains(pollId, ":city:") {
			pollId = pollId + ":city:" + cityKey
		}
	}
	// Rango: defaults Top-10
	start := int64(0)
	stop := int64(9)
	if s := strings.TrimSpace(c.Query("start")); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil && v >= 0 {
			start = v
		}
	}
	if s := strings.TrimSpace(c.Query("stop")); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil && v >= start {
			stop = v
		}
	}
	items, err := h.agg.GetLeaderboard(c.Request.Context(), pollId, start, stop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	type entry struct {
		Member string `json:"member"`
		Score  int64  `json:"score"`
	}
	out := make([]entry, 0, len(items))
	for _, it := range items {
		out = append(out, entry{Member: it.Member, Score: it.Score})
	}
	c.JSON(http.StatusOK, out)
}

// GetStats maneja GET /api/public/stats/:poll_id
func (h *PublicHandlers) GetStats(c *gin.Context) {
	if h.agg == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service Unavailable", "message": "Aggregates not configured"})
		return
	}
	pollId := strings.TrimSpace(c.Param("poll_id"))
	if pollId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": badRequest, "message": "poll_id requerido"})
		return
	}
	if city := strings.TrimSpace(c.Query("city")); city != "" {
		cityKey := normalizeCityKey(city)
		if !strings.Contains(pollId, ":city:") {
			pollId = pollId + ":city:" + cityKey
		}
	}
	st, err := h.agg.GetStats(c.Request.Context(), pollId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total":   st.Total,
		"last_ts": st.LastTS,
		"version": st.Version,
		"uniques": st.Uniques,
	})
}

// GetCount maneja GET /api/public/count/:poll_id/:member (ZSCORE)
func (h *PublicHandlers) GetCount(c *gin.Context) {
	if h.agg == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service Unavailable", "message": "Aggregates not configured"})
		return
	}
	pollId := strings.TrimSpace(c.Param("poll_id"))
	member := strings.TrimSpace(c.Param("member"))
	if pollId == "" || member == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": badRequest, "message": "parametros requeridos"})
		return
	}
	if city := strings.TrimSpace(c.Query("city")); city != "" {
		cityKey := normalizeCityKey(city)
		if !strings.Contains(pollId, ":city:") {
			pollId = pollId + ":city:" + cityKey
		}
	}
	score, err := h.agg.GetScore(c.Request.Context(), pollId, member)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"member": member, "score": score})
}
