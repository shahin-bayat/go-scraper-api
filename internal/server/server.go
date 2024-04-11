package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
	"github.com/shahin-bayat/scraper-api/internal/config"
	"github.com/shahin-bayat/scraper-api/internal/services"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

func Create(db *sqlx.DB, redis *redis.Client, appConfig *config.AppConfig) (*http.Server, error) {
	store := store.New(db, redis)
	services, err := services.NewServices(appConfig)
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", appConfig.Port),
		Handler:      RegisterRoutes(store, services, appConfig),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}, nil

}
