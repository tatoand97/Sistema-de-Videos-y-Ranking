package consumer

import (
	"context"
	"encoding/json"
	"log/slog"

	"admincache/internal/application/cache"
	"admincache/internal/application/keys"
	"admincache/internal/infrastructure"
	amqp "github.com/rabbitmq/amqp091-go"
)

type voteEvent struct {
	Type           string `json:"type"`
	VideoID        int64  `json:"video_id"`
	UserID         int64  `json:"user_id"`
	City           string `json:"city"`
	IdempotencyKey string `json:"idempotency_key"`
}

func StartVoteConsumer(ch *amqp.Channel, cacheImpl *infrastructure.Cache, cfg infrastructure.Config, log *slog.Logger, stop <-chan struct{}) {
	msgs, err := ch.Consume(cfg.VoteQueue, "admincache-consumer", false, false, false, false, nil)
	if err != nil {
		log.Error("consume failed", "err", err)
		return
	}
	ic := cache.New(cacheImpl)

	go func() {
		for {
			select {
			case <-stop:
				log.Info("vote consumer stopped")
				return
			case m, ok := <-msgs:
				if !ok { return }
				var ev voteEvent
				if err := json.Unmarshal(m.Body, &ev); err != nil {
					log.Error("bad vote event", "err", err)
					m.Nack(false, false)
					continue
				}
				citySlug := keys.SlugCity(ev.City)
				if err := ic.DeleteWildcard(context.Background(), "rank:global:*"); err != nil {
					log.Warn("invalidate global failed", "err", err)
				}
				if citySlug != "" {
					if err := ic.DeleteWildcard(context.Background(), "rank:city:"+citySlug+":*"); err != nil {
						log.Warn("invalidate city failed", "city", citySlug, "err", err)
					}
				}
				m.Ack(false)
			}
		}
	}()
}
