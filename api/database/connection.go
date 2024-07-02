package database

import (
	"context"

	sqlc "github.com/BrunoQuaresma/openticket/api/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	queries *sqlc.Queries
	pgConn  *pgxpool.Pool
}

func Connect(connStr string) (Connection, error) {
	dbCtx := context.Background()
	pgConn, err := pgxpool.New(dbCtx, connStr)
	if err != nil {
		return Connection{}, err
	}
	return Connection{
		pgConn:  pgConn,
		queries: sqlc.New(pgConn),
	}, nil
}

func (db *Connection) Close() {
	db.pgConn.Close()
}

func (db *Connection) Queries() *sqlc.Queries {
	return db.queries
}

type txFn func(ctx context.Context, qtx *sqlc.Queries, tx pgx.Tx) error

func (db *Connection) TX(fn txFn) error {
	ctx := context.Background()
	tx, err := db.pgConn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = fn(ctx, db.Queries().WithTx(tx), tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
