# database/sql

A minimal SQL database abstraction for Go with OpenTelemetry instrumentation, designed to support multiple database drivers (PostgreSQL, MySQL, SQL Server, SQLite) without changing application code.

## Usage

```go
import "github.com/stonear/go-dev-toolkit/database/sql"

func main() {
	// use PostgreSQL
	db, err := sql.NewPostgres(
		sql.WithUsername("user"),
		sql.WithPassword("pass"),
		sql.WithHost("localhost"),
		sql.WithPort(5432),
		sql.WithDatabase("mydb"),
	)

	// or use MySQL
	db, err := sql.NewMySQL(
		sql.WithUsername("user"),
		sql.WithPassword("pass"),
		sql.WithHost("localhost"),
		sql.WithPort(3306),
		sql.WithDatabase("mydb"),
	)

	// or use SQL Server
	db, err := sql.NewSQLServer(
		sql.WithUsername("sa"),
		sql.WithPassword("pass"),
		sql.WithHost("localhost"),
		sql.WithPort(1433),
		sql.WithDatabase("mydb"),
		sql.WithInstance("SQLExpress"), // optional named instance
	)

	// or use SQLite
	db, err := sql.NewSQLite(
		sql.WithDatabase("app.db"), // file path or ":memory:"
	)

	// all drivers support timezone (default: "Local")
	db, err := sql.NewPostgres(
		sql.WithUsername("user"),
		sql.WithPassword("pass"),
		sql.WithHost("localhost"),
		sql.WithPort(5432),
		sql.WithDatabase("mydb"),
		sql.WithTimezone("Asia/Jakarta"),
	)
}
```
