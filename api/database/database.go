package database

import (
	"context"

	sqlc "github.com/BrunoQuaresma/openticket/api/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	queries *sqlc.Queries
	conn    *pgxpool.Pool
}

func New(connStr string) (Database, error) {
	dbCtx := context.Background()
	dbConn, err := pgxpool.New(dbCtx, connStr)
	if err != nil {
		return Database{}, err
	}
	return Database{
		conn:    dbConn,
		queries: sqlc.New(dbConn),
	}, nil
}

func (db *Database) Close() {
	db.conn.Close()
}

func (db *Database) Queries() *sqlc.Queries {
	return db.queries
}

type txFn func(ctx context.Context, qtx *sqlc.Queries, tx pgx.Tx) error

func (db *Database) TX(fn txFn) error {
	ctx := context.Background()
	tx, err := db.conn.Begin(ctx)
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
