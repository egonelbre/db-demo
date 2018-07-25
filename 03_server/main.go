package main

import (
	"log"
	"net/http"
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

	log.Println("Started listening on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
