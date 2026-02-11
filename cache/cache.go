package cache

import (
	"context"
	"encoding/json"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
}

type Option func(*Config)

type Config struct {
	Host     string
	Port     int
	Password string
	Database int // Redis/Valkey DB number
}

func Set[T any](ctx context.Context, c Cache, key string, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.Set(ctx, key, data, ttl)
}

func Get[T any](ctx context.Context, c Cache, key string) (T, error) {
	var zero T

	data, err := c.Get(ctx, key)
	if err != nil {
		return zero, err
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return zero, err
	}

	return result, nil
}

func Remember[T any](ctx context.Context, c Cache, key string, ttl time.Duration, fn func() (T, error)) (T, error) {
	var zero T

	data, err := c.Get(ctx, key)
	if err == nil {
		var result T
		if err := json.Unmarshal(data, &result); err != nil {
			return zero, err
		}
		return result, nil
	}

	result, err := fn()
	if err != nil {
		return zero, err
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		return zero, err
	}

	if err := c.Set(ctx, key, bytes, ttl); err != nil {
		return zero, err
	}

	return result, nil
}

func WithHost(host string) Option {
	return func(c *Config) {
		c.Host = host
	}
}

func WithPort(port int) Option {
	return func(c *Config) {
		c.Port = port
	}
}

func WithPassword(password string) Option {
	return func(c *Config) {
		c.Password = password
	}
}

func WithDatabase(database int) Option {
	return func(c *Config) {
		c.Database = database
	}
}
