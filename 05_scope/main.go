package main

import (
	"context"
	"log"
	"net/http"

	"github.com/egonelbre/db-demo/05_scope/pgdb"
	"github.com/egonelbre/db-demo/05_scope/site"
)

func main() {
	ctx := context.Background()

	db, err := pgdb.New(ctx, "user=dbdemo password=dbdemo dbname=dbdemo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := site.NewServer(db)

	log.Println("Started listening on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
