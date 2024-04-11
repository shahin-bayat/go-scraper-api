package store

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/shahin-bayat/scraper-api/internal/config"
)

func NewPostgresStore(cfg *config.PostgresConfig) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.PgUser, cfg.PgPassword, cfg.PgHost, cfg.PgPort, cfg.PgDatabase)
	if cfg.PgInternalUrl != "" {
		connStr = cfg.PgInternalUrl
	}
	//INFO: to fix the issue with railway app : https://docs.railway.app/guides/private-networking#initialization-time
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
