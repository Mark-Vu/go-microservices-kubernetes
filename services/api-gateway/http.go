package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"ride-sharing/shared/contracts"
	"ride-sharing/shared/httputil"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the incoming request from the user/frontend
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

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		httputil.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "internal encoding error"})
		return
	}

	// 4. Create a context with a timeout (Best practice for inter-service calls)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 5. Build the request to the internal trip-service
	targetURL := "http://trip-service:8083/preview"
	outgoingReq, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		httputil.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to create internal request"})
		return
	}
	outgoingReq.Header.Set("Content-Type", "application/json")

	// 6. Execute the call
	client := &http.Client{}
	resp, err := client.Do(outgoingReq)
	if err != nil {
		httputil.WriteJson(w, http.StatusServiceUnavailable, map[string]string{"error": "trip-service is unreachable"})
		return
	}

	// 7. IMPORTANT: Close the body of the response from trip-service
	defer resp.Body.Close()

	var tripResult any
	if err := json.NewDecoder(resp.Body).Decode(&tripResult); err != nil {
		httputil.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to parse response"})
		return
	}

	response := contracts.APIResponse{Data: tripResult}
	httputil.WriteJson(w, http.StatusOK, response)
}
