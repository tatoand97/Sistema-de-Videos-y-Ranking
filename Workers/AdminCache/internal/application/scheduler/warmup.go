package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	mathrand "math/rand"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"admincache/internal/application/keys"
	"admincache/internal/application/ranking"
	"admincache/internal/infrastructure"
)

const (
	scopeGlobal = "global"
	scopeCity   = "city"
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

type cityIndexEntry struct {
	SchemaVersion string    `json:"schema_version"`
	UpdatedAt     time.Time `json:"updated_at"`
	Cities        []string  `json:"cities"`
}

type cycleStats struct {
	refreshed         int
	skippedLock       int
	fetchErrors       int
	validationErrors  int
	writeErrors       int
	staleServed       int
	staleExpired      int
	lockErrors        int
	lockReleaseErrors int
}

type scopeInput struct {
	scope    string
	cityName string
	citySlug string
}

func StartWarmup(comp ranking.Computer, cache *infrastructure.Cache, cfg infrastructure.Config, log *slog.Logger, stop <-chan struct{}) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshIntervalSeconds) * time.Second)
	defer ticker.Stop()

	rand := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

	runCycle := func(trigger string) {
		ctx := context.Background()
		stats := &cycleStats{}
		start := time.Now()

		refreshGlobal(ctx, comp, cache, cfg, log, stats)
		refreshCities(ctx, comp, cache, cfg, log, stats, rand)

		log.Info("warmup cycle completed",
			"trigger", trigger,
			"duration_ms", time.Since(start).Milliseconds(),
			"refreshed", stats.refreshed,
			"skipped_lock", stats.skippedLock,
			"fetch_errors", stats.fetchErrors,
			"validation_errors", stats.validationErrors,
			"write_errors", stats.writeErrors,
			"stale_served", stats.staleServed,
			"stale_expired", stats.staleExpired,
			"lock_errors", stats.lockErrors,
			"lock_release_errors", stats.lockReleaseErrors,
		)
	}

	runCycle("startup")

	for {
		select {
		case <-stop:
			log.Info("warmup scheduler stopped")
			return
		case <-ticker.C:
			runCycle("interval")
		}
	}
}

func refreshGlobal(ctx context.Context, comp ranking.Computer, cache *infrastructure.Cache, cfg infrastructure.Config, log *slog.Logger, stats *cycleStats) {
	processScope(ctx, comp, cache, cfg, log, stats, scopeInput{scope: scopeGlobal})
}

func refreshCities(ctx context.Context, comp ranking.Computer, cache *infrastructure.Cache, cfg infrastructure.Config, log *slog.Logger, stats *cycleStats, rand *mathrand.Rand) {
	if len(cfg.WarmCities) == 0 {
		return
	}

	slugs := make([]string, 0, len(cfg.WarmCities))
	for idx, name := range cfg.WarmCities {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		slug := keys.SlugCity(trimmed)
		slugs = append(slugs, slug)

		processScope(ctx, comp, cache, cfg, log, stats, scopeInput{
			scope:    scopeCity,
			cityName: trimmed,
			citySlug: slug,
		})

		if cfg.CityBatchSize > 0 && (idx+1)%cfg.CityBatchSize == 0 {
			delay := time.Duration(150+rand.Intn(350)) * time.Millisecond
			time.Sleep(delay)
		}
	}

	indexPayload := cityIndexEntry{
		SchemaVersion: cfg.SchemaVersion,
		UpdatedAt:     time.Now().UTC(),
		Cities:        slugs,
	}
	if data, err := json.Marshal(indexPayload); err == nil {
		if err := cache.SetBytes(ctx, keys.CityIndex(cfg.SchemaVersion), data); err != nil {
			log.Warn("failed to update city index", "err", err)
		}
	} else {
		log.Warn("failed to marshal city index", "err", err)
	}
}

func processScope(ctx context.Context, comp ranking.Computer, cache *infrastructure.Cache, cfg infrastructure.Config, log *slog.Logger, stats *cycleStats, scope scopeInput) {
	lockKey := keys.RankLockGlobal(cfg.SchemaVersion)
	dataKey := keys.RankGlobal(cfg.SchemaVersion)
	description := "global"

	if scope.scope == scopeCity {
		lockKey = keys.RankLockCity(scope.citySlug, cfg.SchemaVersion)
		dataKey = keys.RankCity(scope.citySlug, cfg.SchemaVersion)
		description = fmt.Sprintf("city:%s", scope.citySlug)
	}

	token, acquired, err := cache.AcquireLock(ctx, lockKey)
	if err != nil {
		stats.lockErrors++
		log.Error("failed to acquire cache lock", "scope", description, "err", err)
		return
	}
	if !acquired {
		stats.skippedLock++
		log.Info("lock busy, skipping scope", "scope", description)
		return
	}
	defer func() {
		if err := cache.ReleaseLock(ctx, lockKey, token); err != nil {
			stats.lockReleaseErrors++
			log.Warn("failed to release cache lock", "scope", description, "err", err)
		}
	}()

	items, err := fetchRankingWithRetry(ctx, comp, cfg, scope)
	if err != nil {
		stats.fetchErrors++
		log.Error("ranking fetch failed", "scope", description, "err", err)
		handleStaleAssessment(ctx, cache, dataKey, description, log, stats)
		return
	}

	normalized, err := normalizeRanking(items, cfg.MaxTopUsers)
	if err != nil {
		stats.validationErrors++
		log.Error("ranking validation failed", "scope", description, "err", err)
		handleStaleAssessment(ctx, cache, dataKey, description, log, stats)
		return
	}

	now := time.Now().UTC()
	entry := rankingCacheEntry{
		SchemaVersion: cfg.SchemaVersion,
		Scope:         scope.scope,
		City:          scope.cityName,
		CitySlug:      scope.citySlug,
		AsOf:          now,
		FreshUntil:    now.Add(cache.FreshTTL()),
		StaleUntil:    now.Add(cache.FreshTTL() + cache.MaxStale()),
		Items:         normalized,
	}

	payload, err := json.Marshal(entry)
	if err != nil {
		stats.writeErrors++
		log.Error("failed to marshal cache entry", "scope", description, "err", err)
		return
	}

	if err := cache.SetBytes(ctx, dataKey, payload); err != nil {
		stats.writeErrors++
		log.Error("failed to write cache entry", "scope", description, "err", err)
		return
	}

	stats.refreshed++
	log.Info("ranking refreshed", "scope", description, "count", len(normalized), "as_of", entry.AsOf)
}

func fetchRankingWithRetry(ctx context.Context, comp ranking.Computer, cfg infrastructure.Config, scope scopeInput) ([]ranking.RankItem, error) {
	var city *string
	if scope.scope == scopeCity {
		city = &scope.cityName
	}

	size := cfg.MaxTopUsers * 2
	if size < cfg.MaxTopUsers {
		size = cfg.MaxTopUsers
	}

	timeout := time.Duration(cfg.DBReadTimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 3 * time.Second
	}

	var lastErr error
	for attempt := 0; attempt < cfg.DBMaxRetries; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, timeout)
		items, err := comp.Compute(attemptCtx, city, 1, size)
		cancel()
		if err == nil {
			return items, nil
		}
		lastErr = err
		sleep := time.Duration(1<<attempt) * 150 * time.Millisecond
		time.Sleep(sleep)
	}
	if lastErr == nil {
		lastErr = errors.New("unknown fetch error")
	}
	return nil, lastErr
}

func normalizeRanking(items []ranking.RankItem, limit int) ([]rankingCacheItem, error) {
	if limit <= 0 {
		limit = 10
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Votes == items[j].Votes {
			return items[i].UserID < items[j].UserID
		}
		return items[i].Votes > items[j].Votes
	})

	seen := make(map[int64]struct{})
	seenUsernames := make(map[string]struct{})

	result := make([]rankingCacheItem, 0, min(limit, len(items)))
	rank := 0
	for _, it := range items {
		if it.UserID != 0 {
			if _, exists := seen[it.UserID]; exists {
				return nil, fmt.Errorf("duplicate user_id %d", it.UserID)
			}
			seen[it.UserID] = struct{}{}
		} else {
			if _, exists := seenUsernames[it.Username]; exists {
				return nil, fmt.Errorf("duplicate username %s", it.Username)
			}
			seenUsernames[it.Username] = struct{}{}
		}
		rank++
		result = append(result, rankingCacheItem{
			Rank:     rank,
			UserID:   it.UserID,
			Username: it.Username,
			Score:    it.Votes,
		})
		if rank == limit {
			break
		}
	}

	return result, nil
}

func handleStaleAssessment(ctx context.Context, cache *infrastructure.Cache, key, scope string, log *slog.Logger, stats *cycleStats) {
	raw, err := cache.GetBytes(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Warn("no cached data available", "scope", scope)
			return
		}
		log.Warn("unable to inspect stale cache", "scope", scope, "err", err)
		return
	}

	var entry rankingCacheEntry
	if err := json.Unmarshal(raw, &entry); err != nil {
		log.Warn("failed to decode cached entry", "scope", scope, "err", err)
		return
	}

	now := time.Now().UTC()
	if entry.StaleUntil.IsZero() {
		stats.staleServed++
		log.Warn("serving stale cache without deadline", "scope", scope)
		return
	}
	if now.After(entry.StaleUntil) {
		stats.staleExpired++
		log.Error("stale cache expired", "scope", scope, "stale_until", entry.StaleUntil)
		return
	}

	stats.staleServed++
	log.Warn("serving stale cache", "scope", scope, "stale_until", entry.StaleUntil)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
