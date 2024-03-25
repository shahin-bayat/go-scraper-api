package store

import (
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store interface {
	Health() map[string]string
}
