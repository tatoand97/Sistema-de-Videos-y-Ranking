package scheduler

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"admincache/internal/application/keys"
	"admincache/internal/application/ranking"
	"admincache/internal/infrastructure"
)

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func StartWarmup(comp ranking.Computer, cache *infrastructure.Cache, cfg infrastructure.Config, log *slog.Logger, stop <-chan struct{}) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshIntervalSeconds) * time.Second)
	defer ticker.Stop()

	warm := func() {
		ctx := context.Background()
		size := cfg.PageSizeDefault

		// Global pages
		for p := 1; p <= cfg.WarmPages; p++ {
			items, err := comp.Compute(ctx, nil, p, size)
			if err != nil {
				log.Warn("warmup global compute error", "page", p, "err", err)
				continue
			}
			if err := cache.SetBytes(ctx, keys.RankGlobal(p, size), mustJSON(items)); err != nil {
				log.Warn("warmup global set cache error", "page", p, "err", err)
			}
		}

		// City pages
		for _, c := range cfg.WarmCities {
			cslug := c
			for p := 1; p <= cfg.WarmPages; p++ {
				items, err := comp.Compute(ctx, &cslug, p, size)
				if err != nil {
					log.Warn("warmup city compute error", "city", cslug, "page", p, "err", err)
					continue
				}
				if err := cache.SetBytes(ctx, keys.RankCity(cslug, p, size), mustJSON(items)); err != nil {
					log.Warn("warmup city set cache error", "city", cslug, "page", p, "err", err)
				}
			}
		}
		log.Info("warmup done")
	}

	warm()
	for {
		select {
		case <-stop:
			log.Info("warmup scheduler stopped")
			return
		case <-ticker.C:
			warm()
		}
	}
}
