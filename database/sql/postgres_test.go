package sql

import (
	"testing"
)

func TestPostgresDriver(t *testing.T) {
	driver := &PostgresDriver{}
	db := &DB{
		username: "user",
		password: "pass",
		host:     "localhost",
		port:     5432,
		database: "mydb",
		timezone: "UTC",
	}

	expectedDSN := "postgres://user:pass@localhost:5432/mydb?sslmode=disable&timezone=UTC"
	if dsn := driver.DSN(db); dsn != expectedDSN {
		t.Errorf("expected DSN %s, got %s", expectedDSN, dsn)
	}

	if name := driver.Name(); name != "pgx" {
		t.Errorf("expected name pgx, got %s", name)
	}
}
