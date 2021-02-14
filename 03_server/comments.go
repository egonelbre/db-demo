package main

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Comments struct {
	db *pgxpool.Pool
}

func NewComments(ctx context.Context, params string) (*Comments, error) {
	db, err := pgxpool.Connect(ctx, params)
	if err != nil {
		return nil, err
	}
	repo := &Comments{db}
	return repo, repo.init(ctx)
}

func (repo *Comments) Close() error {
	repo.db.Close()
	return nil
}

func (repo *Comments) init(ctx context.Context) error {
	_, err := repo.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS Comments (
			"User"    TEXT,
			"Comment" TEXT
		)
	`)
	return err
}

func (repo *Comments) Add(ctx context.Context, user, comment string) error {
	_, err := repo.db.Exec(ctx, `INSERT INTO Comments ("User", "Comment") VALUES ($1, $2)`, user, comment)
	return err
}

//gistsnip
func (repo *Comments) List(ctx context.Context) ([]Comment, error) {
	rows, err := repo.db.Query(ctx, `SELECT "User", "Comment" FROM Comments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.User, &comment.Text)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
