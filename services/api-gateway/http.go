package main

import (
	"encoding/json"
	"net/http"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody PreviewTripRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&reqBody); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON payload",
		})
		return
	}

	// validation
	if reqBody.UserID == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": "user ID is required",
		})
		return
	}

	response := contracts.APIResponse{Data: "OK"}
	writeJson(w, http.StatusOK, response)

}
