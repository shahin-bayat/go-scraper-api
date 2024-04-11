package handlers

import (
	"database/sql"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/shahin-bayat/scraper-api/internal/utils"
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

	utils.WriteJSON(w, http.StatusOK, healthStatus, nil)

}
