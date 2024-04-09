package handlers

import (
	"github.com/shahin-bayat/scraper-api/internal/services"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

type Handler struct {
	store    store.Store
	services *services.Services
}

func New(store store.Store, services *services.Services) *Handler {
	return &Handler{
		store:    store,
		services: services,
	}
}
