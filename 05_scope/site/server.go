package site

import (
	"context"
	"net/http"
)

//gistsnip:start:db
type DB interface {
	Comments() Comments
}

type Comments interface {
	Add(ctx context.Context, user, comment string) error
	List(ctx context.Context) ([]Comment, error)
}

type Server struct {
	db DB
}

//gistsnip:end:db

func NewServer(db DB) *Server {
	return &Server{
		db: db,
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		server.HandleList(w, r)
	case "/comment":
		server.HandleAddComment(w, r)
	default:
		ShowErrorPage(w, http.StatusNotFound, "Page not found", nil)
	}
}

//gistsnip:start:db
func (server *Server) HandleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		ShowErrorPage(w, http.StatusMethodNotAllowed, "Invalid method", nil)
		return
	}

	comments, err := server.db.Comments().List(ctx)
	if err != nil {
		ShowErrorPage(w, http.StatusInternalServerError, "Unable to access DB", err)
		return
	}

	ShowCommentsPage(w, comments)
}

//gistsnip:end:db

func (server *Server) HandleAddComment(w http.ResponseWriter, r *http.Request) {
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

	err := server.db.Comments().Add(ctx, user, comment)
	if err != nil {
		ShowErrorPage(w, http.StatusInternalServerError, "Unable to add data", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
