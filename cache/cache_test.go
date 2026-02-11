package cache

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockCache struct {
	data map[string][]byte
	err  error
}

func (m *mockCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if m.err != nil {
		return m.err
	}
	m.data[key] = value
	return nil
}

func (m *mockCache) Get(ctx context.Context, key string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	val, ok := m.data[key]
	if !ok {
		return nil, errors.New("not found")
	}
	return val, nil
}

func (m *mockCache) Del(ctx context.Context, key string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.data, key)
	return nil
}

func TestHelpers(t *testing.T) {
	m := &mockCache{data: make(map[string][]byte)}
	ctx := context.Background()

	type User struct {
		Name string
	}

	user := User{Name: "test"}

	// Test Set
	err := Set(ctx, m, "user:1", user, time.Minute)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test Get
	got, err := Get[User](ctx, m, "user:1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Name != "test" {
		t.Errorf("expected test, got %s", got.Name)
	}

	// Test Get failure (not found)
	_, err = Get[User](ctx, m, "user:2")
	if err == nil {
		t.Fatal("expected error on missing key")
	}

	// Test Del
	err = Del(ctx, m, "user:1")
	if err != nil {
		t.Fatalf("Del failed: %v", err)
	}
	_, err = Get[User](ctx, m, "user:1")
	if err == nil {
		t.Fatal("expected error after Del")
	}

	// Test Set error
	m.err = errors.New("error")
	err = Set(ctx, m, "user:1", user, time.Minute)
	if err == nil {
		t.Fatal("expected error")
	}

	// Test Get error (invalid JSON)
	m.err = nil
	m.data["bad"] = []byte("{invalid")
	_, err = Get[User](ctx, m, "bad")
	if err == nil {
		t.Fatal("expected error on invalid JSON")
	}

	// Test Set error (marshal error)
	err = Set(ctx, m, "fail", make(chan int), time.Minute)
	if err == nil {
		t.Fatal("expected marshal error")
	}

	// Test Del error
	m.err = errors.New("error")
	err = Del(ctx, m, "user:1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRemember(t *testing.T) {
	m := &mockCache{data: make(map[string][]byte)}
	ctx := context.Background()

	type User struct {
		Name string
	}

	called := 0
	fn := func() (User, error) {
		called++
		return User{Name: "db"}, nil
	}

	// 1. Miss - calls fn
	got, err := Remember(ctx, m, "user:1", time.Minute, fn)
	if err != nil {
		t.Fatalf("Remember failed: %v", err)
	}
	if got.Name != "db" {
		t.Errorf("expected db, got %s", got.Name)
	}
	if called != 1 {
		t.Errorf("expected 1 call, got %d", called)
	}

	// 2. Hit - doesn't call fn
	got, err = Remember(ctx, m, "user:1", time.Minute, fn)
	if err != nil {
		t.Fatalf("Remember failed: %v", err)
	}
	if got.Name != "db" {
		t.Errorf("expected db, got %s", got.Name)
	}
	if called != 1 {
		t.Errorf("expected 1 call, got %d", called)
	}

	// 3. Hit - invalid JSON
	m.data["user:1"] = []byte("{invalid")
	_, err = Remember(ctx, m, "user:1", time.Minute, fn)
	if err == nil {
		t.Fatal("expected unmarshal error in Remember hit")
	}

	// 4. Fn error
	fnErr := errors.New("fn error")
	fnWithError := func() (User, error) {
		return User{}, fnErr
	}
	_, err = Remember(ctx, m, "user:2", time.Minute, fnWithError)
	if err != fnErr {
		t.Fatalf("expected %v, got %v", fnErr, err)
	}

	// 5. Marshal error during miss
	type marshaler chan int
	fnChannel := func() (marshaler, error) {
		return make(marshaler), nil
	}
	_, err = Remember(ctx, m, "user:3", time.Minute, fnChannel)
	if err == nil {
		t.Fatal("expected marshal error in Remember miss")
	}

	// 6. Set error during miss
	m.err = errors.New("set error")
	_, err = Remember(ctx, m, "user:4", time.Minute, fn)
	if err == nil {
		t.Fatal("expected set error in Remember miss")
	}
}

func TestOptions(t *testing.T) {
	cfg := &Config{}
	opts := []Option{
		WithHost("host"),
		WithPort(123),
		WithPassword("pass"),
		WithDatabase(1),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Host != "host" {
		t.Errorf("expected host, got %s", cfg.Host)
	}
	if cfg.Port != 123 {
		t.Errorf("expected 123, got %d", cfg.Port)
	}
	if cfg.Password != "pass" {
		t.Errorf("expected pass, got %s", cfg.Password)
	}
	if cfg.Database != 1 {
		t.Errorf("expected 1, got %d", cfg.Database)
	}
}
