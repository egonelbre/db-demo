package pgdb

import (
	"context"

	"github.com/egonelbre/db-demo/05_scope/site"
)

type Comments struct {
	db *DB
}

var _ site.Comments = (*Comments)(nil)

func (repo *Comments) Add(ctx context.Context, user, comment string) error {
	_, err := repo.db.Exec(ctx, `INSERT INTO Comments ("User", "Comment") VALUES ($1, $2)`, user, comment)
	return err
}

func (repo *Comments) List(ctx context.Context) ([]site.Comment, error) {
	rows, err := repo.db.Query(ctx, `SELECT "User", "Comment" FROM Comments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []site.Comment{}
	for rows.Next() {
		var comment site.Comment
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
