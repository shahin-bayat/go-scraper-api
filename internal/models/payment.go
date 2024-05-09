package models

import (
	"strconv"
	"strings"
)

type CreateIntentRequest struct {
	SubscriptionID string `json:"subscription_id"`
}

func (r *CreateIntentRequest) Validate() map[string]string {
	errors := make(map[string]string)
	if r.SubscriptionID == "" {
		errors["subscription_id"] = "subscription_id is required"
	}
	_, err := strconv.Atoi(strings.TrimSpace(r.SubscriptionID))
	if err != nil {
		errors["subscription_id"] = "subscription_id must be a number"
	}
	return errors
}

type CreateIntentResponse struct {
	PaymentIntent string `json:"payment_intent"`
	//EphemeralKey  string `json:"ephemeral_key"`
	Customer string `json:"customer"`
}
