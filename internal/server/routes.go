package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/shahin-bayat/scraper-api/internal/handlers"
	"github.com/shahin-bayat/scraper-api/internal/services"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

func RegisterRoutes(store store.Store, services *services.Services) http.Handler {
	r := chi.NewRouter()
	handlers := handlers.New(store, services)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", handlers.HealthHandler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/auth/{provider}/user-info", handlers.GetUserInfo)
		r.Get("/auth/{provider}/login", handlers.HandleProviderLogin)
		r.Get("/auth/{provider}/callback", handlers.HandleProviderCallback)
		r.Get("/auth/{provider}/logout", handlers.HandleLogout)

		r.Get("/payment/config", handlers.HandlePaymentConfig)
		r.Post("/payment/webhook", handlers.HandlePaymentWebhook)
		r.Post("/payment/intent", handlers.HandlePaymentIntent)

		r.Group(func(r chi.Router) {
			// r.Use(handlers.AuthMiddleware)
			r.Route("/category", func(r chi.Router) {
				r.Get("/", handlers.GetCategories)
				r.Get("/{categoryId}", handlers.GetCategoryDetail)
			})
			r.Get("/question/{questionId}", handlers.GetQuestionDetail)
			r.Get("/image/{filename}", handlers.GetImage)
		})
	})

	return r
}
