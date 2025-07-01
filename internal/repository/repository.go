package repository

import (
	"context"
	"main-service/internal/sl"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
	Ping(context.Context) error
	Close()
}

type Database struct {
	DB     PgxIface
	logger *sl.MyLogger
}

func New(db PgxIface, logger *sl.MyLogger) *Database {
	return &Database{DB: db, logger: logger}
}
