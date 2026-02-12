package log

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

type Slog struct {
	config Config
	logger *slog.Logger
}

func NewSlog(opts ...Option) Log {
	cfg := &Config{
		Level:  LevelInfo,
		Output: os.Stdout,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	s := &Slog{
		config: *cfg,
	}

	s.logger = slog.New(slog.NewMultiHandler(
		slog.NewJSONHandler(s.config.Output, &slog.HandlerOptions{
			Level: toSlogLevel(s.config.Level),
		}),
		otelslog.NewHandler(""),
	))

	return s
}

func (s *Slog) Debug(ctx context.Context, msg string, attrs ...Attr) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, msg, toSlogAttrs(attrs)...)
}

func (s *Slog) Info(ctx context.Context, msg string, attrs ...Attr) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, msg, toSlogAttrs(attrs)...)
}

func (s *Slog) Warn(ctx context.Context, msg string, attrs ...Attr) {
	s.logger.LogAttrs(ctx, slog.LevelWarn, msg, toSlogAttrs(attrs)...)
}

func (s *Slog) Error(ctx context.Context, msg string, attrs ...Attr) {
	s.logger.LogAttrs(ctx, slog.LevelError, msg, toSlogAttrs(attrs)...)
}

func toSlogAttrs(attrs []Attr) []slog.Attr {
	if len(attrs) == 0 {
		return nil
	}

	out := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		out[i] = slog.Any(a.Key, a.Value)
	}
	return out
}

func toSlogLevel(l Level) slog.Level {
	switch l {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
