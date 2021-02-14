package pgdb

import (
	"context"

	"github.com/egonelbre/db-demo/05_scope/site"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func New(ctx context.Context, params string) (*DB, error) {
	db, err := pgxpool.Connect(ctx, params)
	if err != nil {
		return nil, err
	}
	rdb := &DB{db}
	return rdb, rdb.init(ctx)
}

func (db *DB) Close() error {
	db.Pool.Close()
	return nil
}

func (db *DB) init(ctx context.Context) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS Comments (
			"User"    TEXT,
			"Comment" TEXT
		)
	`)
	return err
}

func (db *DB) Comments() site.Comments { return &Comments{db} }
