package sql

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "modernc.org/sqlite"
)

type SQLiteDriver struct {
}

func NewSQLite(opts ...Option) (*sql.DB, error) {
	return New(&SQLiteDriver{}, opts...)
}

func (driver *SQLiteDriver) DSN(db *DB) string {
	return fmt.Sprintf("%s?_loc=%s", db.database, url.QueryEscape(db.timezone))
}

func (driver *SQLiteDriver) Name() string {
	return "sqlite"
}
