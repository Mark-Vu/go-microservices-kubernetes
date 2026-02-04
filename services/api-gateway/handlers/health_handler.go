package handlers

import (
	"net/http"
	"time"

	"ride-sharing/shared/httputil"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

// HandleHealth handles health check requests
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Service:   "api-gateway",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	httputil.WriteJson(w, http.StatusOK, response)
}
