package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/shahin-bayat/scraper-api/internal/ecosystem"

	_ "github.com/joho/godotenv/autoload"
)

func Create(eco ecosystem.Ecosystem) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      RegisterRoutes(eco),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

}
