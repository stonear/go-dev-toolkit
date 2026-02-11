package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache"
)

type MemcachedCache struct {
	client *otelmemcache.Client
}

func NewMemcached(opts ...Option) *MemcachedCache {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	client := otelmemcache.NewClientWithTracing(
		memcache.New(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
	)

	return &MemcachedCache{client: client}
}

func (m *MemcachedCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	return m.client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(ttl.Seconds()),
	})
}

func (m *MemcachedCache) Get(_ context.Context, key string) ([]byte, error) {
	item, err := m.client.Get(key)
	if err != nil {
		return nil, err
	}

	return item.Value, nil
}
