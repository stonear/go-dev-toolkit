package client

import (
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Config defines the configuration for the HTTP client.
type Config struct {
	Timeout               time.Duration
	MaxIdleConns          int
	MaxIdleConnsPerHost   int
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
}

// Option defines a functional option for configuring the HTTP client.
type Option func(*Config)

// New creates a new HTTP client with the given options and OpenTelemetry instrumentation.
func New(opts ...Option) *http.Client {
	cfg := &Config{
		Timeout:               30 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          cfg.MaxIdleConns,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
		IdleConnTimeout:       cfg.IdleConnTimeout,
		TLSHandshakeTimeout:   cfg.TLSHandshakeTimeout,
		ExpectContinueTimeout: cfg.ExpectContinueTimeout,
	}

	client := &http.Client{
		Timeout:   cfg.Timeout,
		Transport: otelhttp.NewTransport(transport),
	}

	return client
}

// WithTimeout sets the timeout for the HTTP client.
func WithTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.Timeout = timeout
	}
}

// WithMaxIdleConns sets the maximum number of idle connections.
func WithMaxIdleConns(maxIdleConns int) Option {
	return func(cfg *Config) {
		cfg.MaxIdleConns = maxIdleConns
	}
}

// WithMaxIdleConnsPerHost sets the maximum number of idle connections per host.
func WithMaxIdleConnsPerHost(maxIdleConnsPerHost int) Option {
	return func(cfg *Config) {
		cfg.MaxIdleConnsPerHost = maxIdleConnsPerHost
	}
}

// WithIdleConnTimeout sets the idle connection timeout.
func WithIdleConnTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.IdleConnTimeout = timeout
	}
}
