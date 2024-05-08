package server

import (
	"github.com/shahin-bayat/scraper-api/internal/utils"
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
	h := handlers.New(store, services, appConfig)
	m := middlewares.NewMiddlewares(services.AuthService, store)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route(
		"/api/v1", func(r chi.Router) {
			r.Get("/health", h.HealthHandler)
			r.Get("/auth/{provider}/login", utils.Make(h.HandleProviderLogin))
			r.Get("/auth/{provider}/callback", utils.Make(h.HandleProviderCallback))
			r.Get("/auth/{provider}/user-info", utils.Make(h.GetUserInfo))
			r.Post("/auth/{provider}/logout", utils.Make(h.HandleLogout))

			r.Group(
				func(r chi.Router) {
					r.Use(m.Auth)
					r.Route(
						"/category", func(r chi.Router) {
							r.Get("/", h.GetCategories)
							r.Get("/{categoryId}", h.GetCategoryDetail)
						},
					)
					r.Route(
						"/question", func(r chi.Router) {
							r.Get("/{questionId}", h.GetQuestionDetail)
							r.Get("/supported-languages", h.GetSupportedLanguages)
							r.Post("/bookmark", h.ToggleBookmark)
							r.Get("/bookmark", h.GetBookmarks)
							//r.Post("/user-answer", h.HandleUserAnswer)
						},
					)
					r.Get("/image/{filename}", h.GetImage)

					r.Route(
						"/subscription", func(r chi.Router) {
							r.Get("/", h.GetSubscriptions)
							r.Get("/{subscriptionId}", h.GetSubscriptionDetail)
						},
					)

					r.Route(
						"/payment", func(r chi.Router) {
							r.Get("/config", h.GetPaymentConfig)
							r.Post("/webhook", h.HandlePaymentWebhook)
							r.Post("/intent", h.CreatePaymentIntent)
						},
					)
				},
			)
		},
	)

	return r
}
