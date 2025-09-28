package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements domain Cache interface using Redis.
type RedisCache struct {
	rdb    *redis.Client
	prefix string
}

// MustRedisClient creates a redis client and panics if it cannot ping.
// Reuses the pattern used in Workers/AdminCache.
func MustRedisClient(addr string) *redis.Client {
	r := redis.NewClient(&redis.Options{Addr: addr})
	if err := r.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return r
}

// NewRedisCache builds a cache with a given client and key prefix.
func NewRedisCache(rdb *redis.Client, prefix string) *RedisCache {
	return &RedisCache{rdb: rdb, prefix: prefix}
}

func (c *RedisCache) key(k string) string { return c.prefix + k }

// GetBytes fetches a key and returns its raw bytes.
func (c *RedisCache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return c.rdb.Get(ctx, c.key(key)).Bytes()
}
