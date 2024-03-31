package handlers

import (
	"net/http"

	"github.com/shahin-bayat/scraper-api/internal/utils"
)

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.store.CategoryRepository().GetCategories()
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusBadRequest, categories)

}
