package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, connString string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalln(err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return pool
}
