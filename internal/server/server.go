package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/shahin-bayat/scraper-api/internal/store"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	db   store.Store
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		db:   store.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
