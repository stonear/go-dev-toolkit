package telemetry

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

// Config defines the configuration for the telemetry package.
type Config struct {
	ServiceName    string
	ServiceVersion string
	Exporter       string // "otlp" or "stdout"
}

// Option defines a functional option for configuring the telemetry package.
type Option func(*Config)

// WithServiceName sets the service name.
func WithServiceName(name string) Option {
	return func(c *Config) {
		c.ServiceName = name
	}
}

// WithServiceVersion sets the service version.
func WithServiceVersion(version string) Option {
	return func(c *Config) {
		c.ServiceVersion = version
	}
}

// WithExporterOTLP sets the exporter to OTLP.
func WithExporterOTLP() Option {
	return func(c *Config) {
		c.Exporter = "otlp"
	}
}

// WithExporterStdout sets the exporter to stdout.
func WithExporterStdout() Option {
	return func(c *Config) {
		c.Exporter = "stdout"
	}
}

// New initializes the OpenTelemetry SDK.
func New(ctx context.Context, opts ...Option) (func(context.Context) error, error) {
	cfg := &Config{
		ServiceName:    "unknown-service",
		ServiceVersion: "0.0.0",
		Exporter:       "otlp",
	}

	for _, opt := range opts {
		opt(cfg)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
		),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithHost(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	var shutdownFuncs []func(context.Context) error

	// Shutdown function that calls all registered shutdown functions.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// Handle shutdown if initialization fails.
	handleErr := func(inErr error) error {
		return errors.Join(inErr, shutdown(ctx))
	}

	// Tracing
	var traceExporter sdktrace.SpanExporter
	if cfg.Exporter == "otlp" {
		traceExporter, err = otlptracegrpc.New(ctx)
	} else {
		traceExporter, err = stdouttrace.New()
	}
	if err != nil {
		return nil, handleErr(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	shutdownFuncs = append(shutdownFuncs, tp.Shutdown)

	// Propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Metrics
	var metricExporter sdkmetric.Exporter
	if cfg.Exporter == "otlp" {
		metricExporter, err = otlpmetricgrpc.New(ctx)
	} else {
		metricExporter, err = stdoutmetric.New()
	}
	if err != nil {
		return nil, handleErr(err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)
	shutdownFuncs = append(shutdownFuncs, mp.Shutdown)

	// Logging
	lp := log.NewLoggerProvider(
		log.WithResource(res),
	)
	global.SetLoggerProvider(lp)
	shutdownFuncs = append(shutdownFuncs, lp.Shutdown)

	return shutdown, nil
}
