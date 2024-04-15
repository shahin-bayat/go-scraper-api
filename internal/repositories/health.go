package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type HealthRepository interface {
	HealthCheck(ctx context.Context) error
}

type healthRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

func NewHealthRepository(db *sqlx.DB, redis *redis.Client) HealthRepository {
	return &healthRepository{
		db:    db,
		redis: redis,
	}
}

func (hr *healthRepository) HealthCheck(ctx context.Context) error {
	if err := hr.db.PingContext(ctx); err != nil {
		return err
	}
	if err := hr.redis.Ping(ctx).Err(); err != nil {
		return err
	}
	return nil
}
