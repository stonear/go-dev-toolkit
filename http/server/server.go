package server

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Config defines the configuration for the HTTP server.
type Config struct {
	Addr              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
	HandlerName       string
}

// Server wraps http.Server to provide additional functionality like graceful shutdown.
type Server struct {
	server *http.Server
}

// Option defines a functional option for configuring the HTTP server.
type Option func(*Config)

// New creates a new HTTP server with the given handler, options, and OpenTelemetry instrumentation.
func New(handler http.Handler, opts ...Option) *Server {
	cfg := &Config{
		Addr:              ":3000",
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
		HandlerName:       "http-server",
	}

	for _, opt := range opts {
		opt(cfg)
	}

	otelHandler := otelhttp.NewHandler(handler, cfg.HandlerName)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           otelHandler,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}

	return &Server{server: srv}
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Addr returns the server's network address.
func (s *Server) Addr() string {
	return s.server.Addr
}

// WithAddr sets the address for the HTTP server.
func WithAddr(addr string) Option {
	return func(cfg *Config) {
		cfg.Addr = addr
	}
}

// WithReadTimeout sets the read timeout for the HTTP server.
func WithReadTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.ReadTimeout = timeout
	}
}

// WithReadHeaderTimeout sets the read header timeout for the HTTP server.
func WithReadHeaderTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.ReadHeaderTimeout = timeout
	}
}

// WithWriteTimeout sets the write timeout for the HTTP server.
func WithWriteTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.WriteTimeout = timeout
	}
}

// WithMaxHeaderBytes sets the maximum number of bytes for the request headers.
func WithMaxHeaderBytes(maxHeaderBytes int) Option {
	return func(cfg *Config) {
		cfg.MaxHeaderBytes = maxHeaderBytes
	}
}

// WithHandlerName sets the name for the instrumentation handler.
func WithHandlerName(name string) Option {
	return func(cfg *Config) {
		cfg.HandlerName = name
	}
}
