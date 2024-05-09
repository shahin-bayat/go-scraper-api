package models

type CreateIntentRequest struct {
	SubscriptionID int `json:"subscription_id"`
}

func (r *CreateIntentRequest) Validate() map[string]string {
	errors := make(map[string]string)
	if r.SubscriptionID == 0 {
		errors["subscription_id"] = "subscription_id is required"
	}
	return errors
}

type CreateIntentResponse struct {
	PaymentIntent string `json:"payment_intent"`
	//EphemeralKey  string `json:"ephemeral_key"`
	Customer string `json:"customer"`
}
