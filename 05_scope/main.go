package main

import (
	"log"
	"net/http"

	"github.com/egonelbre/db-demo/05_scope/pgdb"
	"github.com/egonelbre/db-demo/05_scope/site"

	_ "github.com/lib/pq"
)

func main() {
	db, err := pgdb.New("user=dbdemo password=dbdemo dbname=dbdemo sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	server := site.NewServer(db)

	log.Println("Started listening on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
