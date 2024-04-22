package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/shahin-bayat/scraper-api/internal/utils"
	"net/http"
	"strconv"
)

func (h *Handler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := h.store.SubscriptionRepository().GetSubscriptions()
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, subscriptions, nil)
}

func (h *Handler) GetSubscriptionDetail(w http.ResponseWriter, r *http.Request) {
	subscriptionId := chi.URLParam(r, "subscriptionId")
	if subscriptionId == "" {
		utils.WriteErrorJSON(w, http.StatusBadRequest, h.store.SubscriptionRepository().ErrorMissingSubscriptionId())
		return
	}
	uintSubscriptionId, err := strconv.Atoi(subscriptionId)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	subscription, err := h.store.SubscriptionRepository().GetSubscriptionDetail(uintSubscriptionId)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, subscription, nil)

}
