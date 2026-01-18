package main

import (
	"log"
	"net/http"
	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8083")
)

func main() {
	log.Printf("Starting trip-service at http://localhost%s", httpAddr)

	inmemRepo := repository.NewInmemRepository()
	service := service.NewTripService(inmemRepo)

	mux := http.NewServeMux()

	httpHandler := h.HttpHandler{Service: service}

	mux.HandleFunc("POST /preview", httpHandler.HandleTripPreview)
	mux.HandleFunc("POST /route", httpHandler.HandleGetRoute)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from API Gateway"))
	})

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server error: %v", err)
	}

}
