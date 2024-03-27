package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/shahin-bayat/scraper-api/internal/store"
)

func Create(db *sqlx.DB) *http.Server {

	store := store.New(db)

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      RegisterRoutes(store),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

}
