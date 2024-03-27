package handlers

import "github.com/shahin-bayat/scraper-api/internal/store"

type Handler struct {
	store store.Store
}

func New(store store.Store) *Handler {
	return &Handler{
		store: store,
	}
}
