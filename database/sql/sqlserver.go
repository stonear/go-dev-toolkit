package sql

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/microsoft/go-mssqldb"
)

type SQLServerDriver struct {
}

func NewSQLServer(opts ...Option) (*sql.DB, error) {
	return New(&SQLServerDriver{}, opts...)
}

func (driver *SQLServerDriver) DSN(db *DB) string {
	query := url.Values{}
	query.Set("database", db.database)
	query.Set("timezone", db.timezone)

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(db.username, db.password),
		Host:     fmt.Sprintf("%s:%d", db.host, db.port),
		RawQuery: query.Encode(),
	}

	if db.instance != "" {
		u.Path = db.instance
	}

	return u.String()
}

func (driver *SQLServerDriver) Name() string {
	return "sqlserver"
}
