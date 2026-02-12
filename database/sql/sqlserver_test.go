package sql

import (
	"testing"
)

func TestSQLServerDriver(t *testing.T) {
	driver := &SQLServerDriver{}
	db := &DB{
		username: "sa",
		password: "password123",
		host:     "localhost",
		port:     1433,
		database: "mydb",
		timezone: "UTC",
	}

	// Basic connection string
	expectedDSN := "sqlserver://sa:password123@localhost:1433?database=mydb&timezone=UTC"
	if dsn := driver.DSN(db); dsn != expectedDSN {
		t.Errorf("expected DSN %s, got %s", expectedDSN, dsn)
	}

	// With instance
	db.instance = "SQLEXPRESS"
	expectedDSNWithInstance := "sqlserver://sa:password123@localhost:1433/SQLEXPRESS?database=mydb&timezone=UTC"
	if dsn := driver.DSN(db); dsn != expectedDSNWithInstance {
		t.Errorf("expected DSN %s, got %s", expectedDSNWithInstance, dsn)
	}

	if name := driver.Name(); name != "sqlserver" {
		t.Errorf("expected name sqlserver, got %s", name)
	}
}
