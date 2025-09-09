package infrastructure

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	rdb     *redis.Client
	prefix  string
	ttlSecs int
}

func MustRedis(addr string) *redis.Client {
	r := redis.NewClient(&redis.Options{Addr: addr})
	if err := r.Ping(context.Background()).Err(); err != nil { panic(err) }
	return r
}

func NewCache(rdb *redis.Client, prefix string, ttl int) *Cache {
	return &Cache{rdb: rdb, prefix: prefix, ttlSecs: ttl}
}

func (c *Cache) key(k string) string { return c.prefix + k }

func (c *Cache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return c.rdb.Get(ctx, c.key(key)).Bytes()
}

func (c *Cache) SetBytes(ctx context.Context, key string, val []byte) error {
	return c.rdb.Set(ctx, c.key(key), val, time.Duration(c.ttlSecs)*time.Second).Err()
}

func (c *Cache) DeleteWildcard(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, next, err := c.rdb.Scan(ctx, cursor, c.key(pattern), 1000).Result()
		if err != nil { return err }
		if len(keys) > 0 { _ = c.rdb.Del(ctx, keys...).Err() }
		cursor = next
		if cursor == 0 { break }
	}
	return nil
}
