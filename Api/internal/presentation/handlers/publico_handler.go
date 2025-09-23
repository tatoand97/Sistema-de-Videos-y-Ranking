package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"api/internal/application/useCase"
	"api/internal/domain"
	"api/internal/domain/interfaces"
	domainresponses "api/internal/domain/responses"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
)

// PublicHandlers maneja endpoints publicos relacionados a videos.
type PublicHandlers struct {
	service            *useCase.PublicService
	cache              interfaces.Cache
	cacheSchemaVersion string
}

// NewPublicHandlers mantiene compatibilidad para tests y uso sin cache.
func NewPublicHandlers(service *useCase.PublicService) *PublicHandlers {
	return &PublicHandlers{service: service}
}

// NewPublicHandlersWithCache permite inyectar un cache solo-lectura.
func NewPublicHandlersWithCache(service *useCase.PublicService, cache interfaces.Cache, schemaVersion string) *PublicHandlers {
	return &PublicHandlers{service: service, cache: cache, cacheSchemaVersion: strings.TrimSpace(schemaVersion)}
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

	// 4) Logica de voto via servicio (incluye verificacion de existencia y unicidad)
	if eventIDPtr != nil {
		err = h.service.VotePublicVideoWithEvent(c.Request.Context(), videoID, userID, eventIDPtr)
	} else {
		err = h.service.VotePublicVideo(c.Request.Context(), videoID, userID)
	}
	if err != nil {
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

	if cached, ok, err := h.rankingsFromCache(c.Request.Context(), city); err == nil && ok {
		c.JSON(http.StatusOK, cached)
		return
	}

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

func (h *PublicHandlers) rankingsFromCache(ctx context.Context, city *string) ([]domainresponses.RankingEntry, bool, error) {
	if h.cache == nil {
		return nil, false, nil
	}

	schemaVersion := strings.TrimSpace(h.cacheSchemaVersion)
	if schemaVersion == "" {
		return nil, false, nil
	}

	key := fmt.Sprintf("rank:global:%s", schemaVersion)
	var citySlug string
	if city != nil {
		trimmed := strings.TrimSpace(*city)
		if trimmed != "" {
			citySlug = slugCity(trimmed)
			if citySlug == "" {
				return nil, false, nil
			}
			key = fmt.Sprintf("rank:city:%s:%s", citySlug, schemaVersion)
		}
	}

	raw, err := h.cache.GetBytes(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var entry rankingCacheEntry
	if err := json.Unmarshal(raw, &entry); err != nil {
		return nil, false, err
	}

	if entry.SchemaVersion != "" && entry.SchemaVersion != schemaVersion {
		return nil, false, fmt.Errorf("schema version mismatch: cache=%s expected=%s", entry.SchemaVersion, schemaVersion)
	}

	if citySlug != "" {
		if entry.Scope != cacheScopeCity {
			return nil, false, fmt.Errorf("unexpected scope %s for city cache", entry.Scope)
		}
		if entry.CitySlug != "" && entry.CitySlug != citySlug {
			return nil, false, fmt.Errorf("city slug mismatch: cache=%s expected=%s", entry.CitySlug, citySlug)
		}
	}

	now := time.Now().UTC()
	if entry.StaleUntil.IsZero() || now.After(entry.StaleUntil) {
		return nil, false, fmt.Errorf("cache entry expired")
	}

	userIDs := collectUserIDs(entry.Items)
	var userMeta map[uint]domainresponses.UserBasic
	if len(userIDs) > 0 {
		basics, err := h.service.UserBasicsByIDs(ctx, userIDs)
		if err != nil {
			return nil, false, err
		}
		userMeta = make(map[uint]domainresponses.UserBasic, len(basics))
		for _, ub := range basics {
			userMeta[ub.UserID] = ub
		}
	}

	resp := make([]domainresponses.RankingEntry, 0, len(entry.Items))
	for _, item := range entry.Items {
		if item.Username == "" {
			continue
		}
		cityPtr := cityFromSources(item.UserID, userMeta, entry.City)
		resp = append(resp, domainresponses.RankingEntry{
			Position: item.Rank,
			Username: item.Username,
			City:     cityPtr,
			Votes:    int(item.Score),
		})
	}
	if len(resp) == 0 {
		return nil, false, fmt.Errorf("cache entry without items")
	}
	return resp, true, nil
}

func collectUserIDs(items []rankingCacheItem) []uint {
	seen := make(map[uint]struct{})
	ids := make([]uint, 0, len(items))
	for _, item := range items {
		if item.UserID <= 0 {
			continue
		}
		id := uint(item.UserID)
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
}

func cityFromSources(userID int64, basics map[uint]domainresponses.UserBasic, fallback string) *string {
	if userID > 0 && basics != nil {
		if ub, ok := basics[uint(userID)]; ok {
			return cloneStringPtr(ub.City)
		}
	}
	return cloneStringValue(fallback)
}

func cloneStringPtr(in *string) *string {
	if in == nil {
		return nil
	}
	out := *in
	return &out
}

func cloneStringValue(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	out := v
	return &out
}

const (
	cacheScopeGlobal = "global"
	cacheScopeCity   = "city"
)

type rankingCacheItem struct {
	Rank     int    `json:"rank"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Score    int64  `json:"score"`
}

type rankingCacheEntry struct {
	SchemaVersion string             `json:"schema_version"`
	Scope         string             `json:"scope"`
	City          string             `json:"city,omitempty"`
	CitySlug      string             `json:"city_slug,omitempty"`
	AsOf          time.Time          `json:"as_of"`
	FreshUntil    time.Time          `json:"fresh_until"`
	StaleUntil    time.Time          `json:"stale_until"`
	Items         []rankingCacheItem `json:"items"`
}

func slugCity(city string) string {
	s := strings.TrimSpace(strings.ToLower(city))
	if s == "" {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))

	for _, ch := range s {
		ch = normalizeRune(ch)
		switch {
		case ch >= 'a' && ch <= 'z':
			b.WriteRune(ch)
		case ch >= '0' && ch <= '9':
			b.WriteRune(ch)
		case ch == '-' || ch == '_':
			b.WriteRune(ch)
		case ch == ' ':
			b.WriteRune('-')
		case unicode.IsSpace(ch):
			b.WriteRune('-')
		}
	}

	return b.String()
}

func normalizeRune(r rune) rune {
	switch r {
	case '\u00e1', '\u00e0', '\u00e4', '\u00e2', '\u00e3', '\u00e5':
		return 'a'
	case '\u00e9', '\u00e8', '\u00eb', '\u00ea':
		return 'e'
	case '\u00ed', '\u00ec', '\u00ef', '\u00ee':
		return 'i'
	case '\u00f3', '\u00f2', '\u00f6', '\u00f4', '\u00f5':
		return 'o'
	case '\u00fa', '\u00f9', '\u00fc', '\u00fb':
		return 'u'
	case '\u00f1':
		return 'n'
	case '\u00e7':
		return 'c'
	case '\u00df':
		return 's'
	default:
		return r
	}
}
