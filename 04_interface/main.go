package main

import (
	"context"
	"log"
	"net/http"

	"github.com/egonelbre/db-demo/04_interface/pgdb"
	"github.com/egonelbre/db-demo/04_interface/site"
)

func main() {
	ctx := context.Background()

	comments, err := pgdb.NewComments(ctx, "user=dbdemo password=dbdemo dbname=dbdemo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer comments.Close()

	server := site.NewServer(comments)

	log.Println("Started listening on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
