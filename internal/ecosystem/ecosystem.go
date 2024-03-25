package ecosystem

import (
	"github.com/jmoiron/sqlx"
	"github.com/shahin-bayat/scraper-api/internal/repositories"
)

type Ecosystem interface {
	DB() *sqlx.DB
	UserRepository() *repositories.UserRepository
}

type ecosystem struct {
	db             *sqlx.DB
	userRepository *repositories.UserRepository
}

var current Ecosystem

func Require() Ecosystem {
	if current != nil {
		return current
	}
	current = &ecosystem{}
	return current
}

// DB Singleton
func (e *ecosystem) DB() *sqlx.DB {
	if e.db == nil {
		e.db = e.requireDB()
	}
	return e.db
}

func (e *ecosystem) UserRepository() *repositories.UserRepository {
	if e.userRepository == nil {
		e.userRepository = repositories.NewUserRepository(e.DB())
	}
	return e.userRepository
}
