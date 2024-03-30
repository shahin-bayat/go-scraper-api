package server

import (
	"net/http"
	"time"

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
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", handlers.HealthHandler)
	r.Get("/auth/profile", handlers.HandleAuthStatus)
	r.Get("/auth/google/login", handlers.HandleLogin)
	r.Get("/auth/google/callback", handlers.HandleCallback)
	r.Get("/auth/logout", handlers.HandleLogout)

	return r
}
