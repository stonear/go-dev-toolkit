package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"
)

// MockDriver implements the Driver interface for testing
type MockDriver struct {
	dsn  string
	name string
}

func (d *MockDriver) DSN(db *DB) string {
	return d.dsn
}

func (d *MockDriver) Name() string {
	if d.name != "" {
		return d.name
	}
	return "mock"
}

// stdMockDriver implements database/sql/driver.Driver
type stdMockDriver struct {
	failPing bool
}

func (d *stdMockDriver) Open(name string) (driver.Conn, error) {
	if name == "fail-open" {
		return nil, errors.New("open failed")
	}
	return &mockConn{failPing: d.failPing}, nil
}

type mockConn struct {
	failPing bool
}

func (c *mockConn) Prepare(query string) (driver.Stmt, error) {
	return &mockStmt{}, nil
}

func (c *mockConn) Ping(ctx context.Context) error {
	if c.failPing {
		return errors.New("ping failed")
	}
	return nil
}

func (c *mockConn) Close() error              { return nil }
func (c *mockConn) Begin() (driver.Tx, error) { return nil, nil }

type mockStmt struct{}

func (s *mockStmt) Close() error                                    { return nil }
func (s *mockStmt) NumInput() int                                   { return 0 }
func (s *mockStmt) Exec(args []driver.Value) (driver.Result, error) { return nil, nil }
func (s *mockStmt) Query(args []driver.Value) (driver.Rows, error)  { return nil, nil }

var mockDrv = &stdMockDriver{}

func init() {
	sql.Register("mock", mockDrv)
}

func TestOptions(t *testing.T) {
	db := &DB{}
	opts := []Option{
		WithUsername("user"),
		WithPassword("pass"),
		WithHost("host"),
		WithPort(1234),
		WithDatabase("db"),
		WithMaxIdleCount(5),
		WithMaxOpen(10),
		WithMaxLifetime(1 * time.Hour),
		WithMaxIdleTime(30 * time.Minute),
		WithInstance("inst"),
		WithTimezone("Asia/Jakarta"),
	}

	for _, opt := range opts {
		opt(db)
	}

	if db.username != "user" {
		t.Errorf("expected username user, got %s", db.username)
	}
	if db.password != "pass" {
		t.Errorf("expected password pass, got %s", db.password)
	}
	if db.host != "host" {
		t.Errorf("expected host host, got %s", db.host)
	}
	if db.port != 1234 {
		t.Errorf("expected port 1234, got %d", db.port)
	}
	if db.database != "db" {
		t.Errorf("expected database db, got %s", db.database)
	}
	if db.maxIdleCount != 5 {
		t.Errorf("expected maxIdleCount 5, got %d", db.maxIdleCount)
	}
	if db.maxOpen != 10 {
		t.Errorf("expected maxOpen 10, got %d", db.maxOpen)
	}
	if db.maxLifetime != 1*time.Hour {
		t.Errorf("expected maxLifetime 1h, got %v", db.maxLifetime)
	}
	if db.maxIdleTime != 30*time.Minute {
		t.Errorf("expected maxIdleTime 30m, got %v", db.maxIdleTime)
	}
	if db.instance != "inst" {
		t.Errorf("expected instance inst, got %s", db.instance)
	}
	if db.timezone != "Asia/Jakarta" {
		t.Errorf("expected timezone Asia/Jakarta, got %s", db.timezone)
	}
}

func TestNew_Success(t *testing.T) {
	mockDrv.failPing = false
	mock := &MockDriver{name: "mock"}
	db, err := New(mock)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	if db == nil {
		t.Fatal("expected db to be not nil")
	}
	_ = db.Close()
}

func TestNew_OpenFailure(t *testing.T) {
	mock := &MockDriver{name: "not-registered"}
	_, err := New(mock)
	if err == nil {
		t.Fatal("expected error on non-registered driver, got nil")
	}
}

func TestNew_PingFailure(t *testing.T) {
	mockDrv.failPing = true
	mock := &MockDriver{name: "mock"}
	_, err := New(mock)
	if err == nil {
		t.Fatal("expected error on ping failure, got nil")
	}
	mockDrv.failPing = false // reset for other tests
}

func TestConvenienceFunctions(t *testing.T) {
	// These call New(...)
	_, _ = NewPostgres(WithHost("localhost"), WithPort(5432))
	_, _ = NewMySQL(WithHost("localhost"), WithPort(3306))
	_, _ = NewSQLite(WithDatabase(":memory:"))
	_, _ = NewSQLServer(WithHost("localhost"), WithPort(1433))
}

func TestClose(t *testing.T) {
	// Test close on nil connection
	db := &DB{}
	if err := db.Close(); err != nil {
		t.Errorf("expected nil error on nil connection close, got %v", err)
	}
}

func TestDB_Close(t *testing.T) {
	stdDB, _ := sql.Open("mock", "any")
	db := &DB{
		conn: stdDB,
	}
	if err := db.Close(); err != nil {
		t.Errorf("expected nil error on close, got %v", err)
	}
}
