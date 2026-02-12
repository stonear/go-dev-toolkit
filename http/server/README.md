# http/server

A production-ready HTTP server wrapper for Go with OpenTelemetry instrumentation, graceful shutdown support, and security-aware configurations.

## Features

- **OpenTelemetry Instrumentation**: Built-in support for tracing and metrics using `otelhttp`.
- **Graceful Shutdown**: Easy-to-use `Shutdown` method to ensure active requests finish processing.
- **Security-Aware**: Optimized defaults for headers and timeouts to mitigate common attacks like Slowloris.
- **Functional Options**: Flexible configuration for addresses, timeouts, and more.

## Usage

```go
import (
	"context"
	"net/http"
	"time"
	"github.com/stonear/go-dev-toolkit/http/server"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Create a new server
	srv := server.New(handler,
		server.WithAddr(":3000"),
		server.WithReadTimeout(5*time.Second),
		server.WithWriteTimeout(10*time.Second),
	)

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// handle error
		}
	}()

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// handle error
	}
}
```

## Options

- `WithAddr(addr string)`: Sets the server address (default: `:3000`).
- `WithReadTimeout(d time.Duration)`: Sets the maximum duration for reading the entire request.
- `WithReadHeaderTimeout(d time.Duration)`: Sets the maximum duration for reading the request headers.
- `WithWriteTimeout(d time.Duration)`: Sets the maximum duration for writing the response.
- `WithMaxHeaderBytes(n int)`: Sets the maximum number of bytes the server will read parsing the request header's keys and values.
- `WithHandlerName(name string)`: Sets the name used for OpenTelemetry instrumentation.
