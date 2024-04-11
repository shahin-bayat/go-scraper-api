package handlers

import (
	"fmt"
	"net/http"

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
	var payload models.CreateIntentRequest
	if err := utils.DecodeRequestBody(r, &payload); err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	stripe.Key = h.appConfig.StripeSecretKey

	// TODO: with subscription_id fetch subscription details and fill the params like Amount, Currency, Customer, PaymentMethod
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(1099),
		Currency: stripe.String(string(stripe.CurrencyEUR)),
		// TODO: payment methods

	}

	pi, err := paymentintent.New(params)
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			utils.WriteErrorJSON(w, stripeErr.HTTPStatusCode, stripeErr)
			return
		} else {
			utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
	}

	resp := models.CreateIntentResponse{
		ClientSecret: pi.ClientSecret,
	}

	utils.WriteJSON(w, http.StatusOK, resp, nil)
}
