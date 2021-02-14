package main

import (
	"context"
	"log"
	"net/http"
)

type Comment struct {
	User string
	Text string
}

func main() {
	ctx := context.Background()

	//gistsnip:start:main
	comments, err := NewComments(ctx, "user=dbdemo password=dbdemo dbname=dbdemo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer comments.Close()

	server := NewServer(comments)

	log.Println("Started listening on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
	//gistsnip:end:main
}
