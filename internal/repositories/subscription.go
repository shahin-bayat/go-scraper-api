package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/shahin-bayat/scraper-api/internal/models"
)

type SubscriptionRepository interface {
	GetSubscriptions() ([]models.Subscription, error)
	GetSubscriptionDetail(subscriptionId uint) (models.Subscription, error)
	GetSubscriptionById(subscriptionId uint) (models.Subscription, error)
}

type subscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) SubscriptionRepository {
	return &subscriptionRepository{
		db: db,
	}
}

func (sr *subscriptionRepository) GetSubscriptions() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := sr.db.Select(&subscriptions, "SELECT * FROM subscriptions ORDER BY id"); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (sr *subscriptionRepository) GetSubscriptionDetail(subscriptionId uint) (models.Subscription, error) {
	var subscription models.Subscription
	if err := sr.db.Get(&subscription, "SELECT * FROM subscriptions WHERE id = $1", subscriptionId); err != nil {
		return subscription, err
	}
	return subscription, nil
}

func (sr *subscriptionRepository) GetSubscriptionById(subscriptionId uint) (models.Subscription, error) {
	var subscription models.Subscription
	if err := sr.db.Get(&subscription, "SELECT * FROM subscriptions WHERE id = $1", subscriptionId); err != nil {
		return subscription, err
	}
	return subscription, nil
}
