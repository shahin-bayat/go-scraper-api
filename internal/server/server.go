package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
	"github.com/shahin-bayat/scraper-api/internal/config"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

func Create(db *sqlx.DB, redis *redis.Client) (*http.Server, error) {
	store := store.New(db, redis)
	config, err := config.New()
	if err != nil {
		return nil, err
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      RegisterRoutes(store, config),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}, nil

}
