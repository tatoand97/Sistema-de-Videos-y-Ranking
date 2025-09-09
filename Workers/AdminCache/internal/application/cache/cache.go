package cache

import (
	"context"
	"encoding/json"

	"admincache/internal/infrastructure"
)

type ICache interface {
	GetBytes(ctx context.Context, key string) ([]byte, error)
	SetBytes(ctx context.Context, key string, val []byte) error
	DeleteWildcard(ctx context.Context, pattern string) error
}

type Impl struct{ inner *infrastructure.Cache }

func New(c *infrastructure.Cache) *Impl { return &Impl{inner: c} }

func (i *Impl) GetBytes(ctx context.Context, key string) ([]byte, error) { return i.inner.GetBytes(ctx, key) }
func (i *Impl) SetBytes(ctx context.Context, key string, val []byte) error { return i.inner.SetBytes(ctx, key, val) }
func (i *Impl) DeleteWildcard(ctx context.Context, pattern string) error { return i.inner.DeleteWildcard(ctx, pattern) }

func (i *Impl) SetJSON(ctx context.Context, key string, v any) error {
	b, err := json.Marshal(v)
	if err != nil { return err }
	return i.SetBytes(ctx, key, b)
}
