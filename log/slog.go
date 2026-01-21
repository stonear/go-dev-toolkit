package log

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type Slog struct {
	level  Level
	writer io.Writer

	logger *slog.Logger
}

type SlogOption func(*Slog)

func NewSlog(opts ...SlogOption) Log {
	s := &Slog{
		writer: os.Stdout,
		level:  LevelInfo,
	}

	for _, opt := range opts {
		opt(s)
	}

	// TODO: Add otelslog once Go 1.26 is released, using slog.NewMultiHandler
	s.logger = slog.New(
		slog.NewJSONHandler(s.writer, &slog.HandlerOptions{
			Level: toSlogLevel(s.level),
		}),
	)

	return s
}

func WithSlogLevel(level Level) SlogOption {
	return func(s *Slog) {
		s.level = level
	}
}

func WithSlogOutput(w io.Writer) SlogOption {
	return func(s *Slog) {
		s.writer = w
	}
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
