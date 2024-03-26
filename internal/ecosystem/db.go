package ecosystem

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func (e *ecosystem) requireDB() *sqlx.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	if os.Getenv("DB_INTERNAL_URL") != "" {
		connStr = os.Getenv("DB_INTERNAL_URL")
	}
	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}
	return db
}
