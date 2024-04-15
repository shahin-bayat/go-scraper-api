package middlewares

import (
	"github.com/shahin-bayat/scraper-api/internal/services"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

type Middlewares struct {
	authService services.AuthService
	store       store.Store
}

func NewMiddlewares(authService services.AuthService, store store.Store) *Middlewares {
	return &Middlewares{
		authService: authService,
		store:       store,
	}
}
