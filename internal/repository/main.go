package repository

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"log/slog"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Close()
}

type Database struct {
	DB     PgxIface
	logger *slog.Logger
}

func New(db PgxIface, logger *slog.Logger) *Database {
	return &Database{DB: db, logger: logger}
}
