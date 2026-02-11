package log

import (
	"context"
	"io"
)

type Log interface {
	Debug(ctx context.Context, msg string, attrs ...Attr)
	Info(ctx context.Context, msg string, attrs ...Attr)
	Warn(ctx context.Context, msg string, attrs ...Attr)
	Error(ctx context.Context, msg string, attrs ...Attr)
}

type Attr struct {
	Key   string
	Value any
}

func Any(key string, value any) Attr {
	return Attr{
		Key:   key,
		Value: value,
	}
}

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Config struct {
	Level  Level
	Output io.Writer
}

type Option func(*Config)

func WithLevel(level Level) Option {
	return func(c *Config) {
		c.Level = level
	}
}

func WithOutput(output io.Writer) Option {
	return func(c *Config) {
		c.Output = output
	}
}
