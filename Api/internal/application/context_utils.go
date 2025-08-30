package application

import (
	"context"
	"time"
)

const defaultTimeout = 3 * time.Second

func CtxTO(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, hasDeadline := ctx.Deadline(); hasDeadline {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, defaultTimeout)
}
