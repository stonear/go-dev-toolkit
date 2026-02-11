package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeyotel"
)

type ValkeyCache struct {
	client valkey.Client
}

func NewValkey(opts ...Option) (*ValkeyCache, error) {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	client, err := valkeyotel.NewClient(valkey.ClientOption{
		InitAddress: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Password:    cfg.Password,
		SelectDB:    cfg.Database,
	})
	if err != nil {
		return nil, err
	}

	return &ValkeyCache{client: client}, nil
}

func (v *ValkeyCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	cmd := v.client.B().Set().Key(key).Value(string(value)).Ex(ttl).Build()
	return v.client.Do(ctx, cmd).Error()
}

func (v *ValkeyCache) Get(ctx context.Context, key string) ([]byte, error) {
	cmd := v.client.B().Get().Key(key).Build()
	return v.client.Do(ctx, cmd).AsBytes()
}

func (v *ValkeyCache) Del(ctx context.Context, key string) error {
	cmd := v.client.B().Del().Key(key).Build()
	return v.client.Do(ctx, cmd).Error()
}
