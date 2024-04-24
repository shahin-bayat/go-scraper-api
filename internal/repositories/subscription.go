package repositories

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/shahin-bayat/scraper-api/internal/models"
)

var (
	ErrorMissingSubscriptionId = errors.New("subscriptionRepository id is required")
	ErrorGetSubscriptions      = errors.New("error getting subscriptions")
	ErrorGetSubscriptionDetail = errors.New("error getting subscription detail")
)

type SubscriptionRepository interface {
	GetSubscriptions() ([]models.Subscription, error)
	GetSubscriptionDetail(subscriptionId int) (models.Subscription, error)
	GetSubscriptionById(subscriptionId int) (models.Subscription, error)
	ErrorMissingSubscriptionId() error
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
	err := sr.db.Select(&subscriptions, "SELECT * FROM subscriptions")
	if err != nil {
		return nil, ErrorGetSubscriptions
	}
	return subscriptions, nil
}

func (sr *subscriptionRepository) GetSubscriptionDetail(subscriptionId int) (models.Subscription, error) {
	var subscription models.Subscription
	err := sr.db.Get(&subscription, "SELECT * FROM subscriptions WHERE id = $1", subscriptionId)
	if err != nil {
		return subscription, ErrorGetSubscriptionDetail
	}
	return subscription, nil
}

func (sr *subscriptionRepository) GetSubscriptionById(subscriptionId int) (models.Subscription, error) {
	var subscription models.Subscription
	err := sr.db.Get(&subscription, "SELECT * FROM subscriptions WHERE id = $1", subscriptionId)
	if err != nil {
		return subscription, ErrorGetSubscriptionDetail
	}
	return subscription, nil
}

func (sr *subscriptionRepository) ErrorMissingSubscriptionId() error {
	return ErrorMissingSubscriptionId
}
