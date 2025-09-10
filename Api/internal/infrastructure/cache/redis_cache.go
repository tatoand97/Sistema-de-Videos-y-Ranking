package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements domain Cache interface using Redis.
type RedisCache struct {
	rdb     *redis.Client
	prefix  string
	ttlSecs int
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

// NewRedisCache builds a cache with a given client, key prefix and TTL seconds.
func NewRedisCache(rdb *redis.Client, prefix string, ttl int) *RedisCache {
	return &RedisCache{rdb: rdb, prefix: prefix, ttlSecs: ttl}
}

func (c *RedisCache) key(k string) string { return c.prefix + k }

// GetBytes fetches a key and returns its raw bytes.
func (c *RedisCache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return c.rdb.Get(ctx, c.key(key)).Bytes()
}

// SetBytes stores raw bytes at key with the configured TTL.
func (c *RedisCache) SetBytes(ctx context.Context, key string, val []byte) error {
	return c.rdb.Set(ctx, c.key(key), val, time.Duration(c.ttlSecs)*time.Second).Err()
}

// DeleteWildcard deletes keys matching the provided pattern (with prefix applied).
func (c *RedisCache) DeleteWildcard(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, next, err := c.rdb.Scan(ctx, cursor, c.key(pattern), 1000).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			_ = c.rdb.Del(ctx, keys...).Err()
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}

// SetNX sets a key only if it does not already exist, with the provided TTL.
// Returns true if the key was set, false if it already existed.
func (c *RedisCache) SetNX(ctx context.Context, key string, val []byte, ttl time.Duration) (bool, error) {
	ok, err := c.rdb.SetNX(ctx, c.key(key), val, ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}
