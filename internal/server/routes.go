package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shahin-bayat/scraper-api/internal/ecosystem"
	"github.com/shahin-bayat/scraper-api/internal/handlers"
)

func RegisterRoutes(eco ecosystem.Ecosystem) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", handlers.HealthHandler(eco))

	return r
}
