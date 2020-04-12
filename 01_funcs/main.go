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

func initDatabse(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Comments (
			"User"    TEXT,
			"Comment" TEXT
		)
	`)
	return err
}

func addComment(db *sql.DB, user, comment string) error {
	_, err := db.Exec(`INSERT INTO Comments ("User", "Comment") VALUES ($1, $2)`, user, comment)
	return err
}

//gistsnip:start:list
func listComments(db *sql.DB) ([]Comment, error) {
	rows, err := db.Query(`SELECT "User", "Comment" FROM Comments`)
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

func main() {
	db, err := sql.Open("postgres", "user=dbdemo password=dbdemo dbname=dbdemo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err := initDatabse(db); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			ShowErrorPage(w, http.StatusMethodNotAllowed, "Invalid method", nil)
			return
		}

		comments, err := listComments(db)
		if err != nil {
			ShowErrorPage(w, http.StatusInternalServerError, "Unable to access DB", err)
			return
		}

		ShowCommentsPage(w, comments)
	})
	//gistsnip:end:list

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

		err := addComment(db, user, comment)
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
	//gistsnip:start:list
}

//gistsnip:end:list
