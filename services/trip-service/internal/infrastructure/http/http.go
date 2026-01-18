package http

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/httputil"
	"ride-sharing/shared/types"
)

type HttpHandler struct {
	Service domain.TripService
}

func (s *HttpHandler) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody PreviewTripRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&reqBody); err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	fare := &domain.RideFareModel{
		UserID: "42",
	}
	trip, err := s.Service.CreateTrip(ctx, fare)
	if err != nil {
		log.Println(err)
	}

	httputil.WriteJson(w, http.StatusCreated, trip)
}

func (s *HttpHandler) HandleGetRoute(w http.ResponseWriter, r *http.Request) {
	var reqBody GetRouteRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&reqBody); err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	trip, err := s.Service.GetRoute(ctx, &reqBody.Pickup, &reqBody.Destination)

	if err != nil {
		log.Println(err)
	}

	httputil.WriteJson(w, http.StatusOK, trip)
}

type PreviewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

type GetRouteRequest struct {
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}
