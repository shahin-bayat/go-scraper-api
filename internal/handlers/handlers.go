package handlers

import (
	"github.com/shahin-bayat/scraper-api/internal/config"
	"github.com/shahin-bayat/scraper-api/internal/services"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

type Handler struct {
	store     store.Store
	services  *services.Services
	appConfig *config.AppConfig
}

func New(store store.Store, services *services.Services, appConfig *config.AppConfig) *Handler {
	return &Handler{
		store:     store,
		services:  services,
		appConfig: appConfig,
	}
}
