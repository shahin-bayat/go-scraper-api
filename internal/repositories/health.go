package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type HealthRepository struct {
	db *sqlx.DB
}

func NewHealthRepository(db *sqlx.DB) *HealthRepository {
	return &HealthRepository{
		db: db,
	}
}

func (hr *HealthRepository) HealthCheck(ctx context.Context) error {
	if err := hr.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
