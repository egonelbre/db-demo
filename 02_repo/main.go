package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Comment struct {
	User string
	Text string
}

type Comments struct {
	db *sql.DB
}

func NewComments(params string) (*Comments, error) {
	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, err
	}
	repo := &Comments{db}
	return repo, repo.init()
}

func (repo *Comments) init() error {
	_, err := repo.db.Exec(`
		CREATE TABLE IF NOT EXISTS Comments (
			"User"    TEXT,
			"Comment" TEXT
		)
	`)
	return err
}

func (repo *Comments) Add(user, comment string) error {
	_, err := repo.db.Exec(`INSERT INTO Comments ("User", "Comment") VALUES ($1, $2)`, user, comment)
	return err
}

func (repo *Comments) List() ([]Comment, error) {
	rows, err := repo.db.Query(`SELECT "User", "Comment" FROM Comments`)
	if err != nil {
		return nil, err
	}

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

func main() {
	commentsRepo, err := NewComments("user=dbdemo password=dbdemo dbname=dbdemo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			ShowErrorPage(w, http.StatusMethodNotAllowed, "Invalid method", nil)
			return
		}

		comments, err := commentsRepo.List()
		if err != nil {
			ShowErrorPage(w, http.StatusInternalServerError, "Unable to access DB", err)
			return
		}

		ShowCommentsPage(w, comments)
	})

	http.HandleFunc("/comment", func(w http.ResponseWriter, r *http.Request) {
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

		err := commentsRepo.Add(user, comment)
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
