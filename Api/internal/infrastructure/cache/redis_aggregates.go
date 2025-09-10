package cache

import (
	"context"
	"fmt"
	"strconv"

	"api/internal/domain/interfaces"
	"github.com/redis/go-redis/v9"
)

// RedisAggregates implements interfaces.Aggregates using Redis primitives.
// Keys are prefixed with `prefix` to avoid collisions and may include a tenant.
type RedisAggregates struct {
	rdb    *redis.Client
	prefix string
	tenant string
}

// NewRedisAggregates builds a new aggregates client.
func NewRedisAggregates(rdb *redis.Client, prefix, tenant string) *RedisAggregates {
	return &RedisAggregates{rdb: rdb, prefix: prefix, tenant: tenant}
}

func (a *RedisAggregates) key(k string) string { return a.prefix + k }

// core composes the core part inside the hash-tag braces for cluster co-location.
func (a *RedisAggregates) core(pollId string) string {
	if a.tenant != "" {
		return fmt.Sprintf("t:%s|votes:%s", a.tenant, pollId)
	}
	return fmt.Sprintf("votes:%s", pollId)
}

func (a *RedisAggregates) lbKey(pollId string) string { return fmt.Sprintf("lb:{%s}", a.core(pollId)) }
func (a *RedisAggregates) statsKey(pollId string) string {
	return fmt.Sprintf("stats:{%s}", a.core(pollId))
}
func (a *RedisAggregates) lastTSKey(pollId string) string {
	return fmt.Sprintf("stats:{%s}.last_ts", a.core(pollId))
}
func (a *RedisAggregates) hllKey(pollId string) string {
	return fmt.Sprintf("hll:{%s}", a.core(pollId))
}

// UpdateAfterVote performs ZINCRBY, HINCRBY total, SET last_ts, optional PFADD and HINCRBY version in a pipeline.
func (a *RedisAggregates) UpdateAfterVote(ctx context.Context, pollId, member string, uniqueMember *string, tsUnix int64, incVersion bool) error {
	pipe := a.rdb.TxPipeline()
	pipe.ZIncrBy(ctx, a.key(a.lbKey(pollId)), 1, member)
	pipe.HIncrBy(ctx, a.key(a.statsKey(pollId)), "total", 1)
	pipe.Set(ctx, a.key(a.lastTSKey(pollId)), strconv.FormatInt(tsUnix, 10), 0)
	if uniqueMember != nil {
		pipe.PFAdd(ctx, a.key(a.hllKey(pollId)), *uniqueMember)
	}
	if incVersion {
		pipe.HIncrBy(ctx, a.key(a.statsKey(pollId)), "version", 1)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (a *RedisAggregates) GetLeaderboard(ctx context.Context, pollId string, start, stop int64) ([]interfaces.LeaderboardEntry, error) {
	res, err := a.rdb.ZRevRangeWithScores(ctx, a.key(a.lbKey(pollId)), start, stop).Result()
	if err != nil {
		return nil, err
	}
	out := make([]interfaces.LeaderboardEntry, 0, len(res))
	for _, z := range res {
		member, _ := z.Member.(string)
		out = append(out, interfaces.LeaderboardEntry{Member: member, Score: int64(z.Score)})
	}
	return out, nil
}

func (a *RedisAggregates) GetStats(ctx context.Context, pollId string) (interfaces.Stats, error) {
	var st interfaces.Stats
	h, err := a.rdb.HGetAll(ctx, a.key(a.statsKey(pollId))).Result()
	if err != nil {
		return st, err
	}
	if v, ok := h["total"]; ok {
		if n, e := strconv.ParseInt(v, 10, 64); e == nil {
			st.Total = n
		}
	}
	if v, ok := h["version"]; ok {
		if n, e := strconv.ParseInt(v, 10, 64); e == nil {
			st.Version = n
		}
	}
	if v, err := a.rdb.Get(ctx, a.key(a.lastTSKey(pollId))).Result(); err == nil {
		if n, e := strconv.ParseInt(v, 10, 64); e == nil {
			st.LastTS = n
		}
	}
	if n, err := a.rdb.PFCount(ctx, a.key(a.hllKey(pollId))).Result(); err == nil {
		st.Uniques = n
	}
	return st, nil
}

func (a *RedisAggregates) GetScore(ctx context.Context, pollId, member string) (int64, error) {
	v, err := a.rdb.ZScore(ctx, a.key(a.lbKey(pollId)), member).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}
