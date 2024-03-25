package ecosystem

import "github.com/jmoiron/sqlx"

type Ecosystem interface {
	DB() *sqlx.DB
}

var current Ecosystem

func Require() Ecosystem {
	if current != nil {
		return current
	}
	current = &ecosystem{}
	return current
}

type ecosystem struct {
	db *sqlx.DB
}

// DB Singleton
func (e *ecosystem) DB() *sqlx.DB {
	if e.db == nil {
		e.db = e.requireDB()
	}
	return e.db
}
