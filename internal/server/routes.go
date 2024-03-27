package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/shahin-bayat/scraper-api/internal/handlers"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

func RegisterRoutes(store store.Store) http.Handler {
	r := chi.NewRouter()
	handlers := handlers.New(store)

	r.Use(middleware.Logger)

	r.Get("/health", handlers.HealthHandler)

	return r
}
