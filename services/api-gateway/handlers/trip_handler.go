package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ride-sharing/services/api-gateway/dto"
	grpcclients "ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/httputil"
)

// TripHandler handles trip-related HTTP requests
type TripHandler struct {
	tripClient *grpcclients.TripServiceClient
}

// NewTripHandler creates a new TripHandler with dependencies injected
func NewTripHandler(tripClient *grpcclients.TripServiceClient) *TripHandler {
	return &TripHandler{
		tripClient: tripClient,
	}
}

// HandleTripPreview handles preview trip requests via gRPC
func (h *TripHandler) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request from the user/frontend
	var reqBody dto.PreviewTripRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&reqBody); err != nil {
		httputil.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON payload",
		})
		return
	}

	if reqBody.UserID == "" {
		httputil.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "user ID is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Call Trip service via gRPC
	tripResult, err := h.tripClient.Client.PreviewTrip(ctx, reqBody.ToProto())

	if err != nil {
		log.Printf("PreviewTrip gRPC error: %v", err)
		httputil.WriteJson(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	response := contracts.APIResponse{Data: tripResult}
	httputil.WriteJson(w, http.StatusOK, response)
}

// HandleGetRoute handles route calculation requests via HTTP (legacy)
func (h *TripHandler) HandleGetRoute(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request
	var reqBody dto.GetRouteRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&reqBody); err != nil {
		httputil.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON payload",
		})
		return
	}

	// Marshal the request body
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		httputil.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "internal encoding error"})
		return
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Build the request to the trip-service
	targetURL := "http://trip-service:8083/route"
	outgoingReq, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		httputil.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to create internal request"})
		return
	}
	outgoingReq.Header.Set("Content-Type", "application/json")

	// Execute the call
	client := &http.Client{}
	resp, err := client.Do(outgoingReq)
	if err != nil {
		httputil.WriteJson(w, http.StatusServiceUnavailable, map[string]string{"error": "trip-service is unreachable"})
		return
	}
	defer resp.Body.Close()

	// Parse the response
	var routeResult any
	if err := json.NewDecoder(resp.Body).Decode(&routeResult); err != nil {
		httputil.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to parse response"})
		return
	}

	response := contracts.APIResponse{Data: routeResult}
	httputil.WriteJson(w, http.StatusOK, response)
}
