package cache

import (
	"api/internal/infrastructure/cache"
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisCache(t *testing.T) {
	client := &redis.Client{}
	prefix := "test:"

	redisCache := cache.NewRedisCache(client, prefix)

	assert.NotNil(t, redisCache)
}

func TestRedisCache_GetBytes(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	redisCache := cache.NewRedisCache(client, "test:")

	ctx := context.Background()
	key := "testkey"

	_, err := redisCache.GetBytes(ctx, key)

	// Expected to fail without Redis connection
	assert.Error(t, err)
}

func TestMustRedisClient_PanicOnFailure(t *testing.T) {
	assert.Panics(t, func() {
		cache.MustRedisClient("invalid:6379")
	})
}
