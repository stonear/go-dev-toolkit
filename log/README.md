# log

A minimal logging abstraction for Go with structured JSON output, designed to support multiple logging backends (slog, zerolog, more coming) without changing application code.

## Usage

```go
import "github.com/stonear/go-dev-toolkit/log"

func main() {
	// use log/slog
	logger := log.NewSlog(
		log.WithLevel(log.LevelDebug),
		log.WithOutput(os.Stdout),
	)

	// or use zerolog
	logger := log.NewZerolog(
		log.WithLevel(log.LevelDebug),
		log.WithOutput(os.Stdout),
	)

	ctx := context.Background()

	logger.Info(ctx, "hello world", log.Any("email", "test@example.com"))
}
```
