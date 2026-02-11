package log

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type Zerolog struct {
	config Config
	logger zerolog.Logger
}

func NewZerolog(opts ...Option) Log {
	cfg := &Config{
		Level:  LevelInfo,
		Output: os.Stdout,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	z := &Zerolog{
		config: *cfg,
	}

	z.logger = zerolog.New(z.config.Output).
		Level(toZerologLevel(z.config.Level)).
		With().
		Timestamp().
		Logger()

	return z
}

func (z *Zerolog) Debug(ctx context.Context, msg string, attrs ...Attr) {
	event := z.logger.Debug()
	applyZerologAttrs(event, attrs)
	event.Msg(msg)
}

func (z *Zerolog) Info(ctx context.Context, msg string, attrs ...Attr) {
	event := z.logger.Info()
	applyZerologAttrs(event, attrs)
	event.Msg(msg)
}

func (z *Zerolog) Warn(ctx context.Context, msg string, attrs ...Attr) {
	event := z.logger.Warn()
	applyZerologAttrs(event, attrs)
	event.Msg(msg)
}

func (z *Zerolog) Error(ctx context.Context, msg string, attrs ...Attr) {
	event := z.logger.Error()
	applyZerologAttrs(event, attrs)
	event.Msg(msg)
}

func applyZerologAttrs(e *zerolog.Event, attrs []Attr) {
	for _, a := range attrs {
		e.Interface(a.Key, a.Value)
	}
}

func toZerologLevel(l Level) zerolog.Level {
	switch l {
	case LevelDebug:
		return zerolog.DebugLevel
	case LevelInfo:
		return zerolog.InfoLevel
	case LevelWarn:
		return zerolog.WarnLevel
	case LevelError:
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
