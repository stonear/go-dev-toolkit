package sql

import (
	"testing"
)

func TestMySQLDriver(t *testing.T) {
	driver := &MySQLDriver{}
	db := &DB{
		username: "user",
		password: "pass",
		host:     "localhost",
		port:     3306,
		database: "mydb",
		timezone: "Local",
	}

	expectedDSN := "user:pass@tcp(localhost:3306)/mydb?parseTime=true&loc=Local"
	if dsn := driver.DSN(db); dsn != expectedDSN {
		t.Errorf("expected DSN %s, got %s", expectedDSN, dsn)
	}

	if name := driver.Name(); name != "mysql" {
		t.Errorf("expected name mysql, got %s", name)
	}
}
