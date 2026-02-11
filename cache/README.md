# cache

A minimal cache abstraction for Go with generic type-safe operations and OpenTelemetry instrumentation, designed to support multiple cache backends (Redis, Memcached, Valkey) without changing application code.

## Usage

```go
import "github.com/stonear/go-dev-toolkit/cache"

func main() {
	// use Redis
	c, err := cache.NewRedis(
		cache.WithHost("localhost"),
		cache.WithPort(6379),
		cache.WithPassword("secret"),
		cache.WithDatabase(0),
	)

	// or use Memcached
	c := cache.NewMemcached(
		cache.WithHost("localhost"),
		cache.WithPort(11211),
	)

	// or use Valkey
	c, err := cache.NewValkey(
		cache.WithHost("localhost"),
		cache.WithPort(6379),
	)

	ctx := context.Background()

	// generic type-safe set
	cache.Set(ctx, c, "user:1", User{Name: "John"}, 10*time.Minute)

	// generic type-safe get
	user, err := cache.Get[User](ctx, c, "user:1")

	// generic type-safe remember (get or compute and set)
	user, err := cache.Remember(ctx, c, "user:1", 10*time.Minute, func() (User, error) {
		return fetchUserFromDB(1)
	})
}
```
