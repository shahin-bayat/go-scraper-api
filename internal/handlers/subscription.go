package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/shahin-bayat/scraper-api/internal/models"
	"github.com/shahin-bayat/scraper-api/internal/utils"
	"net/http"
	"strconv"
)

func (h *Handler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := h.store.SubscriptionRepository().GetSubscriptions()
	var response models.GetSubscriptionsResponse
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	response.Subscriptions = subscriptions
	response.Features = []string{"Ad-Free Experience", "Access to All Questions", "Bookmark Favorite Questions", "Practice Failed Questions", "Practice Challenging Questions", "Train with Image-Based Questions"}

	utils.WriteJSON(w, http.StatusOK, response, nil)
}

func (h *Handler) GetSubscriptionDetail(w http.ResponseWriter, r *http.Request) {
	subscriptionId := chi.URLParam(r, "subscriptionId")
	if subscriptionId == "" {
		utils.WriteErrorJSON(w, http.StatusBadRequest, h.store.SubscriptionRepository().ErrorMissingSubscriptionId())
		return
	}
	intSubscriptionId, err := strconv.Atoi(subscriptionId)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	subscription, err := h.store.SubscriptionRepository().GetSubscriptionDetail(uint(intSubscriptionId))
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, subscription, nil)
}
