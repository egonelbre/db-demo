package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Comment struct {
	User string
	Text string
}

func main() {
	ctx := context.Background()

	//gistsnip:start:list
	db, err := pgxpool.Connect(ctx, "user=dbdemo password=dbdemo dbname=dbdemo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS Comments (
			"User"    TEXT,
			"Comment" TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if r.Method != http.MethodGet {
			ShowErrorPage(w, http.StatusMethodNotAllowed, "Invalid method", nil)
			return
		}

		rows, err := db.Query(ctx, `SELECT "User", "Comment" FROM Comments`)
		if err != nil {
			ShowErrorPage(w, http.StatusInternalServerError, "Unable to access DB", err)
			return
		}
		defer rows.Close()

		comments := []Comment{}
		for rows.Next() {
			var comment Comment
			err := rows.Scan(&comment.User, &comment.Text)
			if err != nil {
				ShowErrorPage(w, http.StatusInternalServerError, "Unable to load data", err)
				return
			}
			comments = append(comments, comment)
		}

		if err := rows.Err(); err != nil {
			ShowErrorPage(w, http.StatusInternalServerError, "Failed to load data from DB", err)
			return
		}

		ShowCommentsPage(w, comments)
	})
	//gistsnip:end:list

	http.HandleFunc("/comment", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if r.Method != http.MethodPost {
			ShowErrorPage(w, http.StatusMethodNotAllowed, "Invalid method", nil)
			return
		}

		if err := r.ParseForm(); err != nil {
			ShowErrorPage(w, http.StatusBadRequest, "Unable to parse data", err)
			return
		}

		user := r.Form.Get("user")
		comment := r.Form.Get("comment")

		_, err = db.Exec(ctx, `INSERT INTO Comments ("User", "Comment") VALUES ($1, $2)`, user, comment)
		if err != nil {
			ShowErrorPage(w, http.StatusInternalServerError, "Unable to add data", err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("Started listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
