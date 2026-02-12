package sql

import (
	"testing"
)

func TestSQLiteDriver(t *testing.T) {
	driver := &SQLiteDriver{}
	db := &DB{
		database: "test.db",
		timezone: "UTC",
	}

	expectedDSN := "test.db?_loc=UTC"
	if dsn := driver.DSN(db); dsn != expectedDSN {
		t.Errorf("expected DSN %s, got %s", expectedDSN, dsn)
	}

	if name := driver.Name(); name != "sqlite" {
		t.Errorf("expected name sqlite, got %s", name)
	}
}
