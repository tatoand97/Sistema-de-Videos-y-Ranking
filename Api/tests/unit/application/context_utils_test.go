package application_test

import (
	"api/internal/application"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCtxTO_SetsDeadlineWhenNone(t *testing.T) {
	ctx := context.Background()
	ctx2, cancel := application.CtxTO(ctx)
	defer cancel()

	d, ok := ctx2.Deadline()
	assert.True(t, ok, "expected deadline to be set")
	// Expect roughly 3 seconds timeout
	until := time.Until(d)
	assert.Greater(t, until, time.Duration(0))
	assert.LessOrEqual(t, until, 3500*time.Millisecond) // allow small slack
}

func TestCtxTO_NoOpWhenHasDeadline(t *testing.T) {
	parent, parentCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer parentCancel()

	ctx2, cancel := application.CtxTO(parent)
	defer cancel()

	d1, ok1 := parent.Deadline()
	d2, ok2 := ctx2.Deadline()
	assert.True(t, ok1 && ok2)
	// Should keep the same deadline (within small tolerance)
	assert.WithinDuration(t, d1, d2, 50*time.Millisecond)
}
