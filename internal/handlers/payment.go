package handlers

import (
	"errors"
	"fmt"
	"github.com/shahin-bayat/scraper-api/internal/middlewares"
	"github.com/shahin-bayat/scraper-api/internal/models"
	"github.com/shahin-bayat/scraper-api/internal/utils"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/webhook"
	"net/http"
)

type paymentConfig struct {
	PublishableKey string `json:"publishableKey"`
}

func (h *Handler) GetPaymentConfig(w http.ResponseWriter, r *http.Request) error {
	config := paymentConfig{
		PublishableKey: h.appConfig.StripePublishableKey,
	}

	utils.WriteJSON(w, http.StatusOK, config, nil)
	return nil
}

func (h *Handler) HandlePaymentWebhook(w http.ResponseWriter, r *http.Request) error {
	body, err := utils.ReadBody(r)
	if err != nil {
		return err
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), h.appConfig.StripeWebhookSecret)
	if err != nil {
		return err
	}

	if event.Type == "checkout.session.completed" {
		fmt.Println("Checkout Session completed!")
	}

	utils.WriteJSON(w, http.StatusOK, nil, nil)
	return nil
}

func (h *Handler) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) error {
	userId, err := middlewares.GetUserIdFromContext(r.Context())
	if err != nil {
		return utils.NewAPIError(http.StatusUnauthorized, h.services.AuthService.ErrorUnauthorized())
	}
	user, err := h.store.UserRepository().GetUserById(userId)
	if err != nil {
		return err
	}

	var req models.CreateIntentRequest
	if err := utils.DecodeRequestBody(r, &req); err != nil {
		return utils.InvalidJSON()
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		return utils.InvalidRequestData(validationErrors)
	}

	subscription, err := h.store.SubscriptionRepository().GetSubscriptionById(uint(req.SubscriptionID))
	if err != nil {
		return err
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
				return utils.NewAPIError(stripeErr.HTTPStatusCode, err)
			}
		}
		if err = h.store.UserRepository().UpdateUser(
			user.ID, &models.UpdateUserRequest{
				StripeCustomerID: stripeCustomer.ID,
			},
		); err != nil {
			return err
		}
	} else {
		stripeCustomer, err = customer.Get(user.StripeCustomerID, nil)
		if err != nil {
			var stripeErr *stripe.Error
			if errors.As(err, &stripeErr) {
				return utils.NewAPIError(stripeErr.HTTPStatusCode, err)
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
			return utils.NewAPIError(stripeErr.HTTPStatusCode, err)
		}
	}

	resp := models.CreateIntentResponse{
		PaymentIntent: pi.ClientSecret,
		Customer:      stripeCustomer.ID,
	}

	utils.WriteJSON(w, http.StatusOK, resp, nil)
	return nil
}
