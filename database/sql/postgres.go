package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDriver struct {
}

func NewPostgres(opts ...Option) (*sql.DB, error) {
	return New(&PostgresDriver{}, opts...)
}

func (driver *PostgresDriver) DSN(db *DB) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&timezone=%s", db.username, db.password, db.host, db.port, db.database, db.timezone)
}

func (driver *PostgresDriver) Name() string {
	return "pgx"
}
