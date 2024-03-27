package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/shahin-bayat/scraper-api/internal/repositories"
)

type Store interface {
	UserRepository() *repositories.UserRepository
	HealthRepository() *repositories.HealthRepository
}

type store struct {
	db               *sqlx.DB
	userRepository   *repositories.UserRepository
	healthRepository *repositories.HealthRepository
}

func New(db *sqlx.DB) Store {
	return &store{
		db: db,
	}

}

func (e *store) UserRepository() *repositories.UserRepository {
	if e.userRepository == nil {
		e.userRepository = repositories.NewUserRepository(e.db)
	}
	return e.userRepository
}

func (e *store) HealthRepository() *repositories.HealthRepository {
	if e.healthRepository == nil {
		e.healthRepository = repositories.NewHealthRepository(e.db)
	}
	return e.healthRepository
}
