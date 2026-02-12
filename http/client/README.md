# http/client

A production-ready HTTP client wrapper for Go with OpenTelemetry instrumentation and optimized default configurations.

## Features

- **OpenTelemetry Instrumentation**: Built-in support for tracing and metrics using `otelhttp`.
- **Functional Options**: Easy configuration using functional options.
- **Optimized Defaults**: Performance-tuned default timeouts and connection pooling inspired by industry best practices (e.g., Cloudflare).
- **Security-First**: Configurable timeouts for all stages of the HTTP lifecycle.

## Usage

```go
import (
	"context"
	"time"
	"github.com/stonear/go-dev-toolkit/http/client"
)

func main() {
	// Create a new client with default settings
	c := client.New()

	// Create a client with custom settings
	c = client.New(
		client.WithTimeout(10*time.Second),
		client.WithMaxIdleConns(50),
		client.WithMaxIdleConnsPerHost(5),
		client.WithIdleConnTimeout(60*time.Second),
	)

	// Use it like a standard *http.Client
	resp, err := c.Get("https://example.com")
	// ...
}
```

## Options

- `WithTimeout(d time.Duration)`: Sets the total request timeout.
- `WithMaxIdleConns(n int)`: Sets the maximum number of idle connections.
- `WithMaxIdleConnsPerHost(n int)`: Sets the maximum number of idle connections per host.
- `WithIdleConnTimeout(d time.Duration)`: Sets the timeout according to which idle connections are closed.
