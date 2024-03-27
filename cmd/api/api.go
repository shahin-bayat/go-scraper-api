package main

import (
	"fmt"
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

	server := server.Create(db)

	err = server.ListenAndServe()

	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
