package main

import (
	"log"

	c "github.com/shahin-bayat/scraper-api/internal/config"
	"github.com/shahin-bayat/scraper-api/internal/server"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

func main() {

	db, err := store.NewPostgresStore(c.PostgresConf)
	if err != nil {
		log.Fatalf("Failed to create Postgres store: %v", err)
	}
	defer db.Close()

	redis, err := store.NewRedisStore(c.RedisConf)
	if err != nil {
		log.Fatalf("Failed to create Redis store: %v", err)
	}

	server, err := server.Create(db, redis, c.AppConf)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Fatal(server.ListenAndServe())

}
