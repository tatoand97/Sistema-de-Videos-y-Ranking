package main

import (
	"os"
	"os/signal"
	"syscall"

	"admincache/internal/application/consumer"
	"admincache/internal/application/ranking"
	"admincache/internal/application/scheduler"
	"admincache/internal/infrastructure"
)

func main() {
	cfg := infrastructure.LoadConfig()
	logger := infrastructure.NewLogger()

	rdb := infrastructure.MustRedis(cfg.RedisAddr)
	db := infrastructure.MustPostgres(cfg.PostgresDSN)
	ch, conn := infrastructure.MustRabbit(cfg.RabbitURL)
	defer conn.Close()
	defer ch.Close()

	infrastructure.EnsureTopology(ch, cfg, logger)

	comp := ranking.NewRankComputer(db)
	cache := infrastructure.NewCache(rdb, cfg.CachePrefix, cfg.CacheTTLSeconds)

	stopConsumer := make(chan struct{})
	go consumer.StartVoteConsumer(ch, cache, cfg, logger, stopConsumer)

	stopWarm := make(chan struct{})
	go scheduler.StartWarmup(comp, cache, cfg, logger, stopWarm)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	logger.Info("shutting down AdminCache...")
	close(stopConsumer)
	close(stopWarm)
}
