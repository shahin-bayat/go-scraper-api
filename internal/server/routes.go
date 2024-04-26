package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/shahin-bayat/scraper-api/internal/config"
	"github.com/shahin-bayat/scraper-api/internal/handlers"
	"github.com/shahin-bayat/scraper-api/internal/middlewares"
	"github.com/shahin-bayat/scraper-api/internal/services"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

func RegisterRoutes(store store.Store, services *services.Services, appConfig *config.AppConfig) http.Handler {
	r := chi.NewRouter()
	handlers := handlers.New(store, services, appConfig)
	middlewares := middlewares.NewMiddlewares(services.AuthService, store)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route(
		"/api/v1", func(r chi.Router) {
			r.Get("/health", handlers.HealthHandler)
			r.Get("/auth/{provider}/user-info", handlers.GetUserInfo)
			r.Get("/auth/{provider}/login", handlers.HandleProviderLogin)
			r.Get("/auth/{provider}/callback", handlers.HandleProviderCallback)
			r.Post("/auth/{provider}/logout", handlers.HandleLogout)

			r.Group(
				func(r chi.Router) {
					r.Use(middlewares.Auth)
					r.Route(
						"/category", func(r chi.Router) {
							r.Get("/", handlers.GetCategories)
							r.Get("/{categoryId}", handlers.GetCategoryDetail)
						},
					)
					r.Route(
						"/question", func(r chi.Router) {
							r.Get("/{questionId}", handlers.GetQuestionDetail)
							r.Get("/supported-languages", handlers.GetSupportedLanguages)
							r.Post("/bookmark", handlers.ToggleBookmark)
							//r.Post("/user-answer", handlers.HandleUserAnswer)
						},
					)
					r.Get("/image/{filename}", handlers.GetImage)

					r.Route(
						"/subscription", func(r chi.Router) {
							r.Get("/", handlers.GetSubscriptions)
							r.Get("/{subscriptionId}", handlers.GetSubscriptionDetail)
						},
					)

					r.Route(
						"/payment", func(r chi.Router) {
							r.Get("/config", handlers.GetPaymentConfig)
							r.Post("/webhook", handlers.HandlePaymentWebhook)
							r.Post("/intent", handlers.CreatePaymentIntent)
						},
					)
				},
			)
		},
	)

	return r
}
