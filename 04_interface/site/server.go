package site

import (
	"context"
	"net/http"
)

//gistsnip:start:interface
type Comments interface {
	Add(ctx context.Context, user, comment string) error
	List(ctx context.Context) ([]Comment, error)
}

type Server struct {
	comments Comments
}

func NewServer(comments Comments) *Server {
	return &Server{
		comments: comments,
	}
}

//gistsnip:end:interface

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

func (server *Server) HandleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		ShowErrorPage(w, http.StatusMethodNotAllowed, "Invalid method", nil)
		return
	}

	comments, err := server.comments.List(ctx)
	if err != nil {
		ShowErrorPage(w, http.StatusInternalServerError, "Unable to access DB", err)
		return
	}

	ShowCommentsPage(w, comments)
}

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

	err := server.comments.Add(ctx, user, comment)
	if err != nil {
		ShowErrorPage(w, http.StatusInternalServerError, "Unable to add data", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
