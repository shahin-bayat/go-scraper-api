package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shahin-bayat/scraper-api/internal/ecosystem"
)

func HealthHandler(eco ecosystem.Ecosystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := make(map[string]string)
		resp["message"] = "Server is up and running."

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("error handling JSON marshal. Err: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err = eco.DB().PingContext(ctx)
		if err != nil {
			log.Fatalf(fmt.Sprintf("db down: %v", err))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}
