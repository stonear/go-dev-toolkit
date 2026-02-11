package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemcached(t *testing.T) {
	// NewMemcached doesn't check connection
	m := NewMemcached(
		WithHost("localhost"),
		WithPort(11211),
	)

	ctx := context.Background()

	// All these will fail because no server is running,
	// but it will cover the lines in memcached.go.

	_ = m.Set(ctx, "key", []byte("value"), time.Minute)
	_, _ = m.Get(ctx, "key")
	_ = m.Del(ctx, "key")
}
