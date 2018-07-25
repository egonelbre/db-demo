package main

import (
	"log"
	"net/http"

	"github.com/egonelbre/db-demo/04_interface/pgdb"
	"github.com/egonelbre/db-demo/04_interface/site"

	_ "github.com/lib/pq"
)

func main() {
	comments, err := pgdb.NewComments("user=dbdemo password=dbdemo dbname=dbdemo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	server := site.NewServer(comments)

	log.Println("Started listening on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
