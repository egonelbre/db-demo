package main

import "net/http"

//gistsnip:start:server
type Server struct {
	comments *Comments
}

//gistsnip:end:server

func NewServer(comments *Comments) *Server {
	return &Server{
		comments: comments,
	}
}

//gistsnip:start:server
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

//gistsnip:end:server

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
