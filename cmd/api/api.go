package main

import (
	"log"

	cfg "github.com/shahin-bayat/scraper-api/internal/config"
	"github.com/shahin-bayat/scraper-api/internal/server"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

func main() {

	db, err := store.NewPostgresStore(cfg.Config.Postgres)
	if err != nil {
		log.Fatalf("Failed to create Postgres store: %v", err)
	}
	defer db.Close()

	redis, err := store.NewRedisStore(cfg.Config.Redis)
	if err != nil {
		log.Fatalf("Failed to create Redis store: %v", err)
	}

	server, err := server.Create(db, redis)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Fatal(server.ListenAndServe())

}
