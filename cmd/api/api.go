package main

import (
	"log"

	"github.com/shahin-bayat/scraper-api/internal/server"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

func main() {
	db, err := store.NewPostgresStore()
	if err != nil {
		log.Fatalf("Failed to create Postgres store: %v", err)
	}
	defer db.Close()

	server, err := server.Create(db)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Fatal(server.ListenAndServe())

}
