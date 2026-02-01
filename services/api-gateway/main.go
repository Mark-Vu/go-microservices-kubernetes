package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcclients "ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/services/api-gateway/handlers"
	"ride-sharing/services/api-gateway/middleware"
	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	log.Println("Starting API Gateway")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize gRPC trip service client with retry
	tripClient, err := InitgRPCServiceWithRetry(
		ctx,
		"trip-service",
		2*time.Second,
		func(ctx context.Context) (*grpcclients.TripServiceClient, error) {
			return grpcclients.NewTripServiceClient()
		},
	)

	if err != nil {
		log.Fatalf("Failed to initialize trip-service client: %v", err)
	}
	defer tripClient.Close()

	log.Println("Trip service gRPC client initialized successfully")

	// Create handlers with dependencies
	tripHandler := handlers.NewTripHandler(tripClient)

	mux := http.NewServeMux()

	// Trip endpoints
	mux.HandleFunc("POST /trip/preview", middleware.EnableCORS(tripHandler.HandleTripPreview))
	mux.HandleFunc("POST /trip/start", middleware.EnableCORS(tripHandler.HandleTripStart))
	mux.HandleFunc("POST /trip/route", middleware.EnableCORS(tripHandler.HandleGetRoute))

	// WebSocket endpoints
	mux.HandleFunc("/ws/drivers", handlers.HandleDriversWebsocket)
	mux.HandleFunc("/ws/riders", handlers.HandleRidersWebsocket)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from API Gateway"))
	})

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Error starting the server: %v", err)

	case sig := <-shutdown:
		log.Printf("Server is shutting down due to %v signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop server gracefully: %v", err)
			server.Close()
		}
	}
}
