package handlers

import (
	grpcclients "ride-sharing/services/api-gateway/grpc_clients"
)

// DriverHandler handles trip-related HTTP requests
type DriverHandler struct {
	driverClient *grpcclients.DriverServiceClient
}

// NewTripHandler creates a new TripHandler with dependencies injected
func NewDriverHandler(driverClient *grpcclients.DriverServiceClient) *DriverHandler {
	return &DriverHandler{
		driverClient: driverClient,
	}
}
