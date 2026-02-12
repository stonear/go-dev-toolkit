package server

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	server := New(handler,
		WithAddr(":9090"),
		WithReadTimeout(20*time.Second),
		WithReadHeaderTimeout(5*time.Second),
		WithWriteTimeout(30*time.Second),
		WithMaxHeaderBytes(2<<20),
		WithHandlerName("test-handler"),
	)

	if server.Addr() != ":9090" {
		t.Errorf("expected addr :9090, got %s", server.Addr())
	}

	// Access the underlying server to verify fields
	s := server.server
	if s.ReadTimeout != 20*time.Second {
		t.Errorf("expected ReadTimeout 20s, got %v", s.ReadTimeout)
	}
	if s.ReadHeaderTimeout != 5*time.Second {
		t.Errorf("expected ReadHeaderTimeout 5s, got %v", s.ReadHeaderTimeout)
	}
	if s.WriteTimeout != 30*time.Second {
		t.Errorf("expected WriteTimeout 30s, got %v", s.WriteTimeout)
	}
	if s.MaxHeaderBytes != 2<<20 {
		t.Errorf("expected MaxHeaderBytes %d, got %d", 2<<20, s.MaxHeaderBytes)
	}
}

func TestServer_Shutdown(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	server := New(handler, WithAddr(":0"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Errorf("failed to shutdown server: %v", err)
	}
}

func TestServer_ListenAndServe(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	server := New(handler, WithAddr(":0"))

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Errorf("failed to shutdown server: %v", err)
	}

	if err := <-errCh; err != nil {
		t.Errorf("server exited with error: %v", err)
	}
}
