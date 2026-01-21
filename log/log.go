package log

import "context"

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
