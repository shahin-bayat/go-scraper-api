package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/shahin-bayat/scraper-api/internal/repositories"
)

type Store interface {
	UserRepository() *repositories.UserRepository
	HealthRepository() *repositories.HealthRepository
}

type store struct {
	db               *sqlx.DB
	redis            *redis.Client
	userRepository   *repositories.UserRepository
	healthRepository *repositories.HealthRepository
}

func New(db *sqlx.DB, redis *redis.Client) Store {
	return &store{
		db:    db,
		redis: redis,
	}

}

func (e *store) UserRepository() *repositories.UserRepository {
	if e.userRepository == nil {
		e.userRepository = repositories.NewUserRepository(e.db, e.redis)
	}
	return e.userRepository
}

func (e *store) HealthRepository() *repositories.HealthRepository {
	if e.healthRepository == nil {
		e.healthRepository = repositories.NewHealthRepository(e.db, e.redis)
	}
	return e.healthRepository
}
