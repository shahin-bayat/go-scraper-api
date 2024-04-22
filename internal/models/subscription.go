package models

import (
	"time"
)

type Subscription struct {
	ID          uint       `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Price       uint       `json:"price" db:"price"`
	Currency    string     `json:"currency" db:"currency"`
	Interval    string     `json:"interval" db:"interval"`
	CreatedAt   time.Time  `json:"-" db:"created_at"`
	UpdatedAt   time.Time  `json:"-" db:"updated_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`
}

type SubscriptionResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Currency    string  `json:"currency"`
}
