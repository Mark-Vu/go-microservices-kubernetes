package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	grpcclients "ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/httputil"
	pb "ride-sharing/shared/proto/trip/v1"
)

type Handler struct {
	tripClient *grpcclients.TripServiceClient
}

func NewHandler(tripClient *grpcclients.TripServiceClient) *Handler {
	return &Handler{
		tripClient: tripClient,
	}
}

func (h *Handler) handleTripPreview(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request from the user/frontend
	var reqBody PreviewTripRequest
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
	tripResult, err := h.tripClient.Client.PreviewTrip(ctx, &pb.PreviewTripRequest{
		UserId: reqBody.UserID,
		StartLocation: &pb.Coordinate{
			Latitude:   reqBody.Pickup.Latitude,
			Longtitude: reqBody.Pickup.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:   reqBody.Destination.Latitude,
			Longtitude: reqBody.Destination.Longitude,
		},
	})

	if err != nil {
		log.Printf("PreviewTrip gRPC error: %v", err)
		httputil.WriteJson(w, http.StatusServiceUnavailable, map[string]string{
			"error": "trip-service error",
		})
		return
	}

	response := contracts.APIResponse{Data: tripResult}
	httputil.WriteJson(w, http.StatusOK, response)
}

func (h *Handler) handleGetRoute(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the incoming request
	var reqBody GetRouteRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&reqBody); err != nil {
		httputil.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON payload",
		})
		return
	}

	// 2. Marshal the request body
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		httputil.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "internal encoding error"})
		return
	}

	// 3. Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// 4. Build the request to the trip-service
	targetURL := "http://trip-service:8083/route"
	outgoingReq, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		httputil.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to create internal request"})
		return
	}
	outgoingReq.Header.Set("Content-Type", "application/json")

	// 5. Execute the call
	client := &http.Client{}
	resp, err := client.Do(outgoingReq)
	if err != nil {
		httputil.WriteJson(w, http.StatusServiceUnavailable, map[string]string{"error": "trip-service is unreachable"})
		return
	}
	defer resp.Body.Close()

	// 6. Parse the response
	var routeResult any
	if err := json.NewDecoder(resp.Body).Decode(&routeResult); err != nil {
		httputil.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to parse response"})
		return
	}

	response := contracts.APIResponse{Data: routeResult}
	httputil.WriteJson(w, http.StatusOK, response)
}
