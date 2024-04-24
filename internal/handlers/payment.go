package handlers

import (
	"errors"
	"fmt"
	"github.com/shahin-bayat/scraper-api/internal/middlewares"
	"github.com/stripe/stripe-go/customer"
	"net/http"
	"strconv"
	"strings"

	"github.com/shahin-bayat/scraper-api/internal/models"
	"github.com/shahin-bayat/scraper-api/internal/utils"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/webhook"
)

type paymentConfig struct {
	PublishableKey string `json:"publishableKey"`
}

func (h *Handler) HandlePaymentConfig(w http.ResponseWriter, r *http.Request) {
	config := paymentConfig{
		PublishableKey: h.appConfig.StripePublishableKey,
	}

	utils.WriteJSON(w, http.StatusOK, config, nil)

}

func (h *Handler) HandlePaymentWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ReadBody(r)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), h.appConfig.StripeWebhookSecret)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	if event.Type == "checkout.session.completed" {
		fmt.Println("Checkout Session completed!")
	}

	utils.WriteJSON(w, http.StatusOK, nil, nil)
}

func (h *Handler) HandlePaymentIntent(w http.ResponseWriter, r *http.Request) {
	userId, err := middlewares.GetUserIdFromContext(r.Context())
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusUnauthorized, h.services.AuthService.ErrorUnauthorized())
		return
	}
	user, err := h.store.UserRepository().GetUserById(userId)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusNotFound, err)
		return
	}

	var payload models.CreateIntentRequest
	if err := utils.DecodeRequestBody(r, &payload); err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	intSubscriptionId, err := strconv.Atoi(strings.TrimSpace(payload.SubscriptionID))
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	subscription, err := h.store.SubscriptionRepository().GetSubscriptionById(intSubscriptionId)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusNotFound, err)
		return
	}

	stripe.Key = h.appConfig.StripeSecretKey

	var stripeCustomer *stripe.Customer

	if user.StripeCustomerID == "" {
		cparams := &stripe.CustomerParams{
			Name:  &user.Name,
			Email: &user.Email,
		}
		stripeCustomer, err = customer.New(cparams)
		if err != nil {
			var stripeErr *stripe.Error
			if errors.As(err, &stripeErr) {
				utils.WriteErrorJSON(w, stripeErr.HTTPStatusCode, stripeErr)
				return
			}
		}
		err = h.store.UserRepository().UpdateUser(
			user.ID, &models.UpdateUserRequest{
				StripeCustomerID: stripeCustomer.ID,
			},
		)
		if err != nil {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		stripeCustomer, err = customer.Get(user.StripeCustomerID, nil)
		if err != nil {
			if err != nil {
				var stripeErr *stripe.Error
				if errors.As(err, &stripeErr) {
					utils.WriteErrorJSON(w, stripeErr.HTTPStatusCode, stripeErr)
					return
				}
			}
		}
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(subscription.Price)),
		Currency: stripe.String(subscription.Currency),
		Customer: stripe.String(stripeCustomer.ID),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		var stripeErr *stripe.Error
		if errors.As(err, &stripeErr) {
			utils.WriteErrorJSON(w, stripeErr.HTTPStatusCode, stripeErr)
			return
		}
	}

	resp := models.CreateIntentResponse{
		PaymentIntent: pi.ClientSecret,
		Customer:      stripeCustomer.ID,
	}

	utils.WriteJSON(w, http.StatusOK, resp, nil)
}
