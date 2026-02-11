package log

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestNewZerolog_Defaults(t *testing.T) {
	var buf bytes.Buffer

	logger := NewZerolog(
		WithOutput(&buf),
	)

	logger.Debug(context.Background(), "debug message")
	logger.Info(context.Background(), "info message")

	out := buf.String()

	if strings.Contains(out, "debug message") {
		t.Fatalf("expected debug message to be filtered by default level")
	}

	if !strings.Contains(out, "info message") {
		t.Fatalf("expected info message to be logged")
	}
}

func TestNewZerolog_WithLevelDebug(t *testing.T) {
	var buf bytes.Buffer

	logger := NewZerolog(
		WithOutput(&buf),
		WithLevel(LevelDebug),
	)

	logger.Debug(context.Background(), "debug message")

	if !strings.Contains(buf.String(), "debug message") {
		t.Fatalf("expected debug message to be logged")
	}
}

func TestNewZerolog_WarnLevelFiltering(t *testing.T) {
	var buf bytes.Buffer

	logger := NewZerolog(
		WithOutput(&buf),
		WithLevel(LevelWarn),
	)

	logger.Info(context.Background(), "info")
	logger.Warn(context.Background(), "warn")
	logger.Error(context.Background(), "error")

	out := buf.String()

	if strings.Contains(out, "info") {
		t.Fatalf("info should be filtered")
	}

	if !strings.Contains(out, "warn") {
		t.Fatalf("warn should be logged")
	}

	if !strings.Contains(out, "error") {
		t.Fatalf("error should be logged")
	}
}

func TestZerolog_WithAttributes(t *testing.T) {
	var buf bytes.Buffer

	logger := NewZerolog(
		WithOutput(&buf),
		WithLevel(LevelDebug),
	)

	logger.Info(
		context.Background(),
		"user logged in",
		Any("user_id", 42),
		Any("email", "test@example.com"),
	)

	var decoded map[string]any
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}

	if decoded["message"] != "user logged in" {
		t.Fatalf("message mismatch")
	}

	if decoded["user_id"].(float64) != 42 {
		t.Fatalf("user_id missing or incorrect")
	}

	if decoded["email"] != "test@example.com" {
		t.Fatalf("email missing or incorrect")
	}
}

func TestZerolog_NoAttributes(t *testing.T) {
	var buf bytes.Buffer

	logger := NewZerolog(
		WithOutput(&buf),
	)

	logger.Info(context.Background(), "hello")

	if !strings.Contains(buf.String(), `"message":"hello"`) {
		t.Fatalf("expected message")
	}
}

func TestNewZerolog_ErrorLevelOnly(t *testing.T) {
	var buf bytes.Buffer

	logger := NewZerolog(
		WithOutput(&buf),
		WithLevel(LevelError),
	)

	logger.Warn(context.Background(), "warn")
	logger.Error(context.Background(), "error")

	out := buf.String()

	if strings.Contains(out, `"message":"warn"`) {
		t.Fatalf("warn should be filtered")
	}

	if !strings.Contains(out, `"message":"error"`) {
		t.Fatalf("error should be logged")
	}
}

func TestApplyZerologAttrs(t *testing.T) {
	var buf bytes.Buffer

	z := zerolog.New(&buf)
	event := z.Info()

	attrs := []Attr{
		{Key: "a", Value: 1},
		{Key: "b", Value: "x"},
	}

	applyZerologAttrs(event, attrs)
	event.Msg("test")

	var decoded map[string]any
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if decoded["a"].(float64) != 1 {
		t.Fatalf("attr a missing")
	}
	if decoded["b"] != "x" {
		t.Fatalf("attr b missing")
	}
}

func TestToZerologLevel(t *testing.T) {
	cases := []struct {
		level Level
	}{
		{LevelDebug},
		{LevelInfo},
		{LevelWarn},
		{LevelError},
		{Level(999)}, // default case
	}

	for _, c := range cases {
		_ = toZerologLevel(c.level)
	}
}
