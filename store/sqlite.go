package store

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

type sqliteRepo struct {
	db *sql.DB
}

// NewSQLite opens the database at path and applies the schema.
func NewSQLite(path string) (*sqliteRepo, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.ExecContext(context.Background(), schema); err != nil {
		db.Close()
		return nil, err
	}
	return &sqliteRepo{db: db}, nil
}

func (r *sqliteRepo) Close() error { return r.db.Close() }

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
