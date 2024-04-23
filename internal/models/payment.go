package models

type CreateIntentRequest struct {
	SubscriptionID string `json:"subscription_id"`
}

type CreateIntentResponse struct {
	PaymentIntent string `json:"payment_intent"`
	//EphemeralKey  string `json:"ephemeral_key"`
	Customer string `json:"customer"`
}
