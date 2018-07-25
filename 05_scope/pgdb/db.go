package pgdb

import (
	"database/sql"

	"github.com/egonelbre/db-demo/05_scope/site"
)

type DB struct {
	*sql.DB
}

func New(params string) (*DB, error) {
	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, err
	}
	rdb := &DB{db}
	return rdb, rdb.init()
}

func (db *DB) init() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Comments (
			"User"    TEXT,
			"Comment" TEXT
		)
	`)
	return err
}

func (db *DB) Comments() site.Comments { return &Comments{db} }
