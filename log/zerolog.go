package log

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Zerolog struct {
	level  Level
	writer io.Writer

	logger zerolog.Logger
}

type ZerologOption func(*Zerolog)

func NewZerolog(opts ...ZerologOption) Log {
	z := &Zerolog{
		writer: os.Stdout,
		level:  LevelInfo,
	}

	for _, opt := range opts {
		opt(z)
	}

	z.logger = zerolog.New(z.writer).
		Level(toZerologLevel(z.level)).
		With().
		Timestamp().
		Logger()

	return z
}

func WithZerologLevel(level Level) ZerologOption {
	return func(z *Zerolog) {
		z.level = level
	}
}

func WithZerologOutput(w io.Writer) ZerologOption {
	return func(z *Zerolog) {
		z.writer = w
	}
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
