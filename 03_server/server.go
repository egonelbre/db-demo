package main

import "net/http"

type Server struct {
	comments *Comments
}

func NewServer(comments *Comments) *Server {
	return &Server{
		comments: comments,
	}
}

func (server *Server) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ShowErrorPage(w, http.StatusMethodNotAllowed, "Invalid method", nil)
		return
	}

	comments, err := server.comments.List()
	if err != nil {
		ShowErrorPage(w, http.StatusInternalServerError, "Unable to access DB", err)
		return
	}

	ShowCommentsPage(w, comments)
}

func (server *Server) HandleAddComment(w http.ResponseWriter, r *http.Request) {
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

	err := server.comments.Add(user, comment)
	if err != nil {
		ShowErrorPage(w, http.StatusInternalServerError, "Unable to add data", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
