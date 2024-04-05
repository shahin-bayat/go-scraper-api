package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type HealthRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

func NewHealthRepository(db *sqlx.DB, redis *redis.Client) *HealthRepository {
	return &HealthRepository{
		db:    db,
		redis: redis,
	}
}

func (hr *HealthRepository) HealthCheck(ctx context.Context) error {
	if err := hr.db.PingContext(ctx); err != nil {
		return err
	}
	if err := hr.redis.Ping(ctx).Err(); err != nil {
		return err
	}
	return nil
}
