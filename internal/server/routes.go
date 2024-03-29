package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/shahin-bayat/scraper-api/internal/config"
	"github.com/shahin-bayat/scraper-api/internal/handlers"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

func RegisterRoutes(store store.Store, config *config.Config) http.Handler {
	r := chi.NewRouter()
	handlers := handlers.New(store, config)

	r.Use(middleware.Logger)

	r.Get("/health", handlers.HealthHandler)
	r.Get("/auth/login", handlers.HandleLogin)
	r.Get("/auth/google/callback", handlers.HandleCallback)

	return r
}
