package application

import (
	"context"
	"time"
)

const defaultTimeout = 3 * time.Second

func CtxTO(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, hasDeadline := ctx.Deadline(); hasDeadline {
		// Intentionally return a no-op cancel when the incoming context already
		// has a deadline. Callers expect a cancel function from this helper, but
		// because we are not creating a derived context in this branch, there is
		// nothing to cancel here. The upstream owner of ctx is responsible for
		// its lifecycle and cancellation.
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, defaultTimeout)
}
