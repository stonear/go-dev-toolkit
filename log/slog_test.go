package log

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func decodeLastJSON(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) == 0 {
		t.Fatalf("no log output")
	}

	var out map[string]any
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	return out
}

func TestNewSlog_Defaults(t *testing.T) {
	var buf bytes.Buffer

	logger := NewSlog(
		WithOutput(&buf),
	)

	logger.Debug(context.Background(), "debug")
	logger.Info(context.Background(), "info")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")

	if len(lines) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(lines))
	}

	out := decodeLastJSON(t, &buf)

	if out["msg"] != "info" {
		t.Fatalf("expected info message")
	}
}

func TestNewSlog_WithLevelDebug(t *testing.T) {
	var buf bytes.Buffer

	logger := NewSlog(
		WithOutput(&buf),
		WithLevel(LevelDebug),
	)

	logger.Debug(context.Background(), "debug")

	out := decodeLastJSON(t, &buf)

	if out["msg"] != "debug" {
		t.Fatalf("expected debug message")
	}
}

func TestNewSlog_WarnLevelFiltering(t *testing.T) {
	var buf bytes.Buffer

	logger := NewSlog(
		WithOutput(&buf),
		WithLevel(LevelWarn),
	)

	logger.Info(context.Background(), "info")
	logger.Warn(context.Background(), "warn")
	logger.Error(context.Background(), "error")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(lines))
	}

	last := decodeLastJSON(t, &buf)
	if last["msg"] != "error" {
		t.Fatalf("expected error message")
	}
}

func TestSlog_WithAttributes(t *testing.T) {
	var buf bytes.Buffer

	logger := NewSlog(
		WithOutput(&buf),
		WithLevel(LevelDebug),
	)

	logger.Info(
		context.Background(),
		"user logged in",
		Any("user_id", 42),
		Any("email", "test@example.com"),
	)

	out := decodeLastJSON(t, &buf)

	if out["msg"] != "user logged in" {
		t.Fatalf("message mismatch")
	}

	if out["user_id"].(float64) != 42 {
		t.Fatalf("user_id missing")
	}

	if out["email"] != "test@example.com" {
		t.Fatalf("email missing")
	}
}

func TestSlog_NoAttributes(t *testing.T) {
	var buf bytes.Buffer

	logger := NewSlog(
		WithOutput(&buf),
	)

	logger.Info(context.Background(), "hello")

	out := decodeLastJSON(t, &buf)

	if out["msg"] != "hello" {
		t.Fatalf("expected message")
	}
}

func TestNewSlog_ErrorLevelOnly(t *testing.T) {
	var buf bytes.Buffer

	logger := NewSlog(
		WithOutput(&buf),
		WithLevel(LevelError),
	)

	logger.Warn(context.Background(), "warn")
	logger.Error(context.Background(), "error")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 log entry")
	}

	out := decodeLastJSON(t, &buf)

	if out["msg"] != "error" {
		t.Fatalf("expected error message")
	}
}

func TestToSlogAttrs(t *testing.T) {
	attrs := []Attr{
		{Key: "a", Value: 1},
		{Key: "b", Value: "x"},
	}

	out := toSlogAttrs(attrs)

	if len(out) != 2 {
		t.Fatalf("expected 2 attrs")
	}

	if out[0].Key != "a" || out[1].Key != "b" {
		t.Fatalf("attr keys mismatch")
	}
}

func TestToSlogLevel(t *testing.T) {
	levels := []Level{
		LevelDebug,
		LevelInfo,
		LevelWarn,
		LevelError,
		Level(999),
	}

	for _, l := range levels {
		_ = toSlogLevel(l)
	}
}
