package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Comment struct {
	User string
	Text string
}

func main() {
	comments, err := NewComments("user=dbdemo password=dbdemo dbname=dbdemo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	server := NewServer(comments)

	http.HandleFunc("/", server.HandleList)
	http.HandleFunc("/comment", server.HandleAddComment)

	log.Println("Started listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
