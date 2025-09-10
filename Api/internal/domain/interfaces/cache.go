package interfaces

import (
	"context"
	"time"
)

// Cache abstracts a key-value cache used by the API layer.
// Implementations live in infrastructure (e.g., Redis).
type Cache interface {
	// GetBytes returns the value stored at key or an error if not found.
	GetBytes(ctx context.Context, key string) ([]byte, error)
	// SetBytes stores the value at key with the configured TTL.
	SetBytes(ctx context.Context, key string, val []byte) error
	// DeleteWildcard deletes keys matching the pattern (supports wildcards).
	DeleteWildcard(ctx context.Context, pattern string) error
	// SetNX stores the value only if key does not exist, with provided TTL.
	// Returns true if the key was set (i.e., first time), false if it already existed.
	SetNX(ctx context.Context, key string, val []byte, ttl time.Duration) (bool, error)
}
