package telemetry

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
	}{
		{
			name: "Default",
			opts: nil,
		},
		{
			name: "Stdout Exporter",
			opts: []Option{
				WithExporterStdout(),
				WithServiceName("stdout-service"),
				WithServiceVersion("1.0.0"),
			},
		},
		{
			name: "OTLP Exporter Explicit",
			opts: []Option{
				WithExporterOTLP(),
				WithServiceName("otlp-service"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			shutdown, err := New(ctx, tt.opts...)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}
			if shutdown == nil {
				t.Fatal("shutdown function is nil")
			}

			// Call shutdown to ensure no panic.
			// For OTLP without collector, shutdown will fail/timeout as it tries to flush.
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()

			if err := shutdown(shutdownCtx); err != nil {
				// Only report error if we expected stdout (which shouldn't fail)
				if tt.name == "Stdout Exporter" {
					t.Errorf("shutdown failed: %v", err)
				}
			}
		})
	}
}
