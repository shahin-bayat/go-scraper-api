package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/shahin-bayat/scraper-api/internal/models"
	"github.com/shahin-bayat/scraper-api/internal/utils"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/webhook"
)

type config struct {
	PublishableKey string `json:"publishableKey"`
}

func (h *Handler) HandlePaymentConfig(w http.ResponseWriter, r *http.Request) {
	publishableKey := os.Getenv("STRIPE_PUBLISHABLE_KEY")
	if publishableKey == "" {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("stripe publishable key not set"))
		return
	}
	config := config{
		PublishableKey: publishableKey,
	}

	utils.WriteJSON(w, http.StatusOK, config, nil)

}

func (h *Handler) HandlePaymentWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ReadBody(r)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if stripeWebhookSecret == "" {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("stripe webhook secret not set"))
		return
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), stripeWebhookSecret)
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

	paymentSecretKey := os.Getenv("STRIPE_SECRET_KEY")
	if paymentSecretKey == "" {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("stripe secret key not set"))
		return
	}

	stripe.Key = paymentSecretKey

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
