package handlers

import (
	"github.com/shahin-bayat/scraper-api/internal/config"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

type Handler struct {
	store  store.Store
	config *config.Config
}

func New(store store.Store, config *config.Config) *Handler {
	return &Handler{
		store:  store,
		config: config,
	}
}
