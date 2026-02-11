package sql

import (
	"database/sql"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelsql"
)

// List of supported drivers available at https://go.dev/wiki/SQLDrivers

type DB struct {
	username string
	password string
	host     string
	port     int
	database string
	instance string // SQL Server named instance
	timezone string // timezone for database connections

	maxIdleCount int           // zero means defaultMaxIdleConns; negative means 0
	maxOpen      int           // <= 0 means unlimited
	maxLifetime  time.Duration // maximum amount of time a connection may be reused
	maxIdleTime  time.Duration // maximum amount of time a connection may be idle before being closed

	conn *sql.DB
}

type Driver interface {
	DSN(db *DB) string
	Name() string
}

type Option func(*DB)

func New(driver Driver, opts ...Option) (*sql.DB, error) {
	db := &DB{
		timezone:     "Local",
		maxIdleCount: 10,
		maxOpen:      100,
		maxLifetime:  1 * time.Hour,
		maxIdleTime:  1 * time.Minute,
	}
	for _, opt := range opts {
		opt(db)
	}

	var err error
	db.conn, err = otelsql.Open(driver.Name(), driver.DSN(db))
	if err != nil {
		return nil, err
	}

	db.conn.SetMaxIdleConns(db.maxIdleCount)
	db.conn.SetMaxOpenConns(db.maxOpen)
	db.conn.SetConnMaxLifetime(db.maxLifetime)
	db.conn.SetConnMaxIdleTime(db.maxIdleTime)
	otelsql.ReportDBStatsMetrics(db.conn)

	if err = db.conn.Ping(); err != nil {
		return nil, err
	}

	return db.conn, nil
}

func WithUsername(username string) Option {
	return func(db *DB) {
		db.username = username
	}
}

func WithPassword(password string) Option {
	return func(db *DB) {
		db.password = password
	}
}

func WithHost(host string) Option {
	return func(db *DB) {
		db.host = host
	}
}

func WithPort(port int) Option {
	return func(db *DB) {
		db.port = port
	}
}

func WithDatabase(database string) Option {
	return func(db *DB) {
		db.database = database
	}
}

func WithMaxIdleCount(maxIdleCount int) Option {
	return func(db *DB) {
		db.maxIdleCount = maxIdleCount
	}
}

func WithMaxOpen(maxOpen int) Option {
	return func(db *DB) {
		db.maxOpen = maxOpen
	}
}

func WithMaxLifetime(maxLifetime time.Duration) Option {
	return func(db *DB) {
		db.maxLifetime = maxLifetime
	}
}

func WithMaxIdleTime(maxIdleTime time.Duration) Option {
	return func(db *DB) {
		db.maxIdleTime = maxIdleTime
	}
}

func WithInstance(instance string) Option {
	return func(db *DB) {
		db.instance = instance
	}
}

func WithTimezone(timezone string) Option {
	return func(db *DB) {
		db.timezone = timezone
	}
}

func (db *DB) Close() error {
	if db.conn == nil {
		return nil
	}

	return db.conn.Close()
}
