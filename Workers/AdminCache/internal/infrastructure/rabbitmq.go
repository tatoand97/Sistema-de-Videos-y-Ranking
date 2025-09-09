package infrastructure

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

func MustRabbit(url string) (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial(url)
	if err != nil { panic(err) }
	ch, err := conn.Channel()
	if err != nil { panic(err) }
	return ch, conn
}

func EnsureTopology(ch *amqp.Channel, cfg Config, log *slog.Logger) {
	_ = ch.ExchangeDeclare(cfg.EventsExchange, "topic", true, false, false, false, nil)
	_ = ch.ExchangeDeclare("dlx.exchange", "direct", true, false, false, false, nil)

	argsRetry := amqp.Table{
		"x-dead-letter-exchange":    cfg.EventsExchange,
		"x-dead-letter-routing-key": cfg.VoteRoutingKey,
		"x-message-ttl":             int32(5000),
	}
	_, _ = ch.QueueDeclare("vote_events.retry.5s", true, false, false, false, argsRetry)

	argsMain := amqp.Table{"x-dead-letter-exchange": "dlx.exchange"}
	_, _ = ch.QueueDeclare(cfg.VoteQueue, true, false, false, false, argsMain)
	_, _ = ch.QueueDeclare("vote_events.dlq", true, false, false, false, nil)

	_ = ch.QueueBind(cfg.VoteQueue, cfg.VoteRoutingKey, cfg.EventsExchange, false, nil)
	_ = ch.QueueBind("vote_events.dlq", "vote_events.dlq", "dlx.exchange", false, nil)

	log.Info("topology ensured", "exchange", cfg.EventsExchange, "queue", cfg.VoteQueue)
}
