package client

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	client := New(
		WithTimeout(10*time.Second),
		WithMaxIdleConns(50),
		WithMaxIdleConnsPerHost(5),
		WithIdleConnTimeout(60*time.Second),
	)

	if client.Timeout != 10*time.Second {
		t.Errorf("expected timeout 10s, got %v", client.Timeout)
	}

	if client.Transport == nil {
		t.Fatal("expected transport to be set")
	}
}
