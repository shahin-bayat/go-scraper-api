package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/shahin-bayat/scraper-api/internal/models"
	"github.com/shahin-bayat/scraper-api/internal/utils"
	"net/http"
	"strconv"
)

var (
	ErrorMissingSubscriptionId = errors.New("subscriptionRepository id is required")
)

func (h *Handler) GetSubscriptions(w http.ResponseWriter, r *http.Request) error {
	subscriptions, err := h.store.SubscriptionRepository().GetSubscriptions()
	if err != nil {
		return err
	}
	var response models.GetSubscriptionsResponse
	response.Subscriptions = subscriptions
	response.Features = []string{"Ad-Free Experience", "Access to All Questions", "Bookmark Favorite Questions", "Practice Failed Questions", "Practice Challenging Questions", "Train with Image-Based Questions"}

	utils.WriteJSON(w, http.StatusOK, response, nil)
	return nil
}

func (h *Handler) GetSubscriptionDetail(w http.ResponseWriter, r *http.Request) error {
	subscriptionId := chi.URLParam(r, "subscriptionId")
	if subscriptionId == "" {
		return utils.NewAPIError(
			http.StatusUnprocessableEntity, ErrorMissingSubscriptionId,
		)
	}
	intSubscriptionId, err := strconv.Atoi(subscriptionId)
	if err != nil {
		return utils.NewAPIError(http.StatusUnprocessableEntity, err)
	}
	subscription, err := h.store.SubscriptionRepository().GetSubscriptionDetail(uint(intSubscriptionId))
	if err != nil {
		return err
	}
	utils.WriteJSON(w, http.StatusOK, subscription, nil)
	return nil
}
