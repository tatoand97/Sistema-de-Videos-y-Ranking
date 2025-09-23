package infrastructure

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	mathrand "math/rand"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheSettings struct {
	Prefix        string
	FreshTTL      time.Duration
	MaxStale      time.Duration
	JitterPercent int
	LockLease     time.Duration
}

type Cache struct {
	rdb           *redis.Client
	prefix        string
	freshTTL      time.Duration
	maxStale      time.Duration
	jitterPercent int
	lockLease     time.Duration

	rndMu sync.Mutex
	rnd   *mathrand.Rand
}

var releaseLockScript = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`)

func MustRedis(addr string) *redis.Client {
	r := redis.NewClient(&redis.Options{Addr: addr})
	if err := r.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return r
}

func NewCache(rdb *redis.Client, settings CacheSettings) *Cache {
	if settings.Prefix == "" {
		settings.Prefix = "videorank:"
	}
	c := &Cache{
		rdb:           rdb,
		prefix:        settings.Prefix,
		freshTTL:      settings.FreshTTL,
		maxStale:      settings.MaxStale,
		jitterPercent: settings.JitterPercent,
		lockLease:     settings.LockLease,
		rnd:           mathrand.New(mathrand.NewSource(time.Now().UnixNano())),
	}
	if c.freshTTL <= 0 {
		c.freshTTL = 15 * time.Minute
	}
	if c.lockLease <= 0 {
		c.lockLease = 10 * time.Second
	}
	return c
}

func (c *Cache) key(k string) string {
	return c.prefix + k
}

func (c *Cache) FreshTTL() time.Duration {
	return c.freshTTL
}

func (c *Cache) MaxStale() time.Duration {
	return c.maxStale
}

func (c *Cache) LockLease() time.Duration {
	return c.lockLease
}

func (c *Cache) ttlWithJitter() time.Duration {
	base := c.freshTTL + c.maxStale
	if base <= 0 {
		return 0
	}
	if c.jitterPercent <= 0 {
		return base
	}
	c.rndMu.Lock()
	defer c.rndMu.Unlock()
	variance := float64(base) * float64(c.jitterPercent) / 100.0
	delta := (c.rnd.Float64()*2 - 1) * variance
	ttl := base + time.Duration(delta)
	if ttl <= 0 {
		return base
	}
	return ttl
}

func (c *Cache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return c.rdb.Get(ctx, c.key(key)).Bytes()
}

func (c *Cache) SetBytes(ctx context.Context, key string, val []byte) error {
	ttl := c.ttlWithJitter()
	return c.rdb.Set(ctx, c.key(key), val, ttl).Err()
}

func (c *Cache) DeleteWildcard(ctx context.Context, pattern string) error {
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

func (c *Cache) AcquireLock(ctx context.Context, key string) (string, bool, error) {
	token, err := randomToken()
	if err != nil {
		return "", false, err
	}
	ok, err := c.rdb.SetNX(ctx, c.key(key), token, c.lockLease).Result()
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	return token, true, nil
}

func (c *Cache) ReleaseLock(ctx context.Context, key, token string) error {
	_, err := releaseLockScript.Run(ctx, c.rdb, []string{c.key(key)}, token).Result()
	if err == redis.Nil {
		return nil
	}
	return err
}

func randomToken() (string, error) {
	b := make([]byte, 16)
	if _, err := cryptorand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
