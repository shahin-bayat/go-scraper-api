package store

import (
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	pgDatabase = os.Getenv("PG_DATABASE")
	pgPassword = os.Getenv("PG_PASSWORD")
	pgUser     = os.Getenv("PG_USER")
	pgPort     = os.Getenv("PG_PORT")
	pgHost     = os.Getenv("PG_HOST")
)

func NewPostgresStore() (*sqlx.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgPort, pgDatabase)
	if os.Getenv("PG_INTERNAL_URL") != "" {
		connStr = os.Getenv("PG_INTERNAL_URL")
	}
	// to fix the issue with railway app : https://docs.railway.app/guides/private-networking#initialization-time
	time.Sleep(3 * time.Second)

	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
