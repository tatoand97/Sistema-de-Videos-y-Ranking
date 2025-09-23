package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"admincache/internal/application/ranking"
	"admincache/internal/application/scheduler"
	"admincache/internal/infrastructure"
)

func main() {
	cfg := infrastructure.LoadConfig()
	logger := infrastructure.NewLogger()

	rdb := infrastructure.MustRedis(cfg.RedisAddr)
	db := infrastructure.MustPostgres(cfg.PostgresDSN)

	comp := ranking.NewRankComputer(db)
	cache := infrastructure.NewCache(rdb, infrastructure.CacheSettings{
		Prefix:        cfg.CachePrefix,
		FreshTTL:      time.Duration(cfg.TTLFreshSeconds) * time.Second,
		MaxStale:      time.Duration(cfg.MaxStaleSeconds) * time.Second,
		JitterPercent: cfg.JitterPercent,
		LockLease:     time.Duration(cfg.LockLeaseSeconds) * time.Second,
	})

	stopWarm := make(chan struct{})
	go scheduler.StartWarmup(comp, cache, cfg, logger, stopWarm)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	logger.Info("shutting down AdminCache...")
	close(stopWarm)
}
