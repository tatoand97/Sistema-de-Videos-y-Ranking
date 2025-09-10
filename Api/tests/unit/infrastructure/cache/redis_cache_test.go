package cache

import (
	"api/internal/infrastructure/cache"
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisCache(t *testing.T) {
	client := &redis.Client{}
	prefix := "test:"
	ttl := 300
	
	redisCache := cache.NewRedisCache(client, prefix, ttl)
	
	assert.NotNil(t, redisCache)
}

func TestRedisCache_GetBytes(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	redisCache := cache.NewRedisCache(client, "test:", 300)
	
	ctx := context.Background()
	key := "testkey"
	
	_, err := redisCache.GetBytes(ctx, key)
	
	// Expected to fail without Redis connection
	assert.Error(t, err)
}

func TestRedisCache_SetBytes(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	redisCache := cache.NewRedisCache(client, "test:", 300)
	
	ctx := context.Background()
	key := "testkey"
	data := []byte("test data")
	
	err := redisCache.SetBytes(ctx, key, data)
	
	// Expected to fail without Redis connection
	assert.Error(t, err)
}

func TestRedisCache_SetNX(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	redisCache := cache.NewRedisCache(client, "test:", 300)
	
	ctx := context.Background()
	key := "testkey"
	data := []byte("test data")
	ttl := 5 * time.Minute
	
	_, err := redisCache.SetNX(ctx, key, data, ttl)
	
	// Expected to fail without Redis connection
	assert.Error(t, err)
}

func TestRedisCache_DeleteWildcard(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	redisCache := cache.NewRedisCache(client, "test:", 300)
	
	ctx := context.Background()
	pattern := "test*"
	
	err := redisCache.DeleteWildcard(ctx, pattern)
	
	// Expected to fail without Redis connection
	assert.Error(t, err)
}

func TestMustRedisClient_PanicOnFailure(t *testing.T) {
	assert.Panics(t, func() {
		cache.MustRedisClient("invalid:6379")
	})
}