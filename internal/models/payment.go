package models

type CreateIntentRequest struct {
	SubscriptionID string `json:"subscription_id"`
}

type CreateIntentResponse struct {
	ClientSecret string `json:"client_secret"`
}
