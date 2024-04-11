package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	healthStatus := make(map[string]string)

	// Check database health
	if err := h.store.HealthRepository().HealthCheck(r.Context()); err != nil {
		if err == sql.ErrConnDone {
			healthStatus["database"] = "error: database connection closed"
		} else if err == redis.ErrClosed {
			healthStatus["redis"] = "error: redis connection closed"
		} else {
			healthStatus["db/redis"] = "error: " + err.Error()
		}
	} else {
		healthStatus["database"] = "ok"
		healthStatus["redis"] = "ok"
	}

	// Add more health checks as needed

	// Marshal health status to JSON
	response, err := json.Marshal(healthStatus)
	if err != nil {
		http.Error(w, "Failed to marshal health status", http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
