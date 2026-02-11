package sql

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDriver struct {
}

func NewMySQL(opts ...Option) (*sql.DB, error) {
	return New(&MySQLDriver{}, opts...)
}

func (driver *MySQLDriver) DSN(db *DB) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s", db.username, db.password, db.host, db.port, db.database, url.QueryEscape(db.timezone))
}

func (driver *MySQLDriver) Name() string {
	return "mysql"
}
