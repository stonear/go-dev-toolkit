package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type mockRedis struct {
	redis.Cmdable
	resSet *redis.StatusCmd
	resGet *redis.StringCmd
	resDel *redis.IntCmd
}

func (m *mockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return m.resSet
}

func (m *mockRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	return m.resGet
}

func (m *mockRedis) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return m.resDel
}

func TestRedis(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		m := &mockRedis{
			resSet: redis.NewStatusCmd(ctx),
			resGet: redis.NewStringCmd(ctx),
			resDel: redis.NewIntCmd(ctx),
		}
		m.resGet.SetVal("value")

		r := &RedisCache{client: m}

		if err := r.Set(ctx, "key", []byte("value"), time.Minute); err != nil {
			t.Errorf("Set failed: %v", err)
		}

		val, err := r.Get(ctx, "key")
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		if string(val) != "value" {
			t.Errorf("expected value, got %s", string(val))
		}

		if err := r.Del(ctx, "key"); err != nil {
			t.Errorf("Del failed: %v", err)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		wantErr := errors.New("redis error")
		m := &mockRedis{
			resSet: redis.NewStatusCmd(ctx),
			resGet: redis.NewStringCmd(ctx),
			resDel: redis.NewIntCmd(ctx),
		}
		m.resSet.SetErr(wantErr)
		m.resGet.SetErr(wantErr)
		m.resDel.SetErr(wantErr)

		r := &RedisCache{client: m}

		if err := r.Set(ctx, "key", []byte("value"), time.Minute); err != wantErr {
			t.Errorf("expected %v, got %v", wantErr, err)
		}

		if _, err := r.Get(ctx, "key"); err != wantErr {
			t.Errorf("expected %v, got %v", wantErr, err)
		}

		if err := r.Del(ctx, "key"); err != wantErr {
			t.Errorf("expected %v, got %v", wantErr, err)
		}
	})
}

func TestNewRedis(t *testing.T) {
	// Should fail because no server is running at this default/random port or address
	// but we can test the constructor and option application.
	r, err := NewRedis(
		WithHost("invalid"),
		WithPort(0),
	)
	// NewRedis doesn't actually fail on connection, it only fails on OTel instrumentation
	// which usually succeeds in creating hooks even if the client is not connected yet.
	if err != nil {
		t.Fatalf("NewRedis failed: %v", err)
	}
	if r == nil {
		t.Fatal("expected RedisCache instance")
	}
}
