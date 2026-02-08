package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"syscall"
	"time"

	"google.golang.org/grpc"

	g "ride-sharing/services/trip-service/internal/infrastructure/grpc"
)

var (
	httpAddr    = env.GetString("HTTP_ADDR", ":8082")
	rabbitmqURI = env.GetString("RABBITMQ_URI", "amqp://guest:guest@localhost:5672")
)

func main() {
	// Start listening for incoming gRPC requests
	lis, err := net.Listen("tcp", httpAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	// Create RabbitMQ client
	rabbitmq, err := messaging.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		log.Fatalf("failed to create RabbitMQ client: %v", err)
	}
	defer rabbitmq.Close()
	log.Printf("RabbitMQ client created successfully")

	inmemRepo := repository.NewInmemRepository()
	svc := service.NewTripService(inmemRepo)

	// Starting grpc server
	grpcServer := grpc.NewServer()
	g.NewGRPCHandler(grpcServer, svc)

	log.Printf("Starting gRPC trip-service on port %s", lis.Addr().String())
	serverError := make(chan error, 1)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			serverError <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("received signal %s: starting graceful shutdown", sig)
	case err := <-serverError:
		// If Serve returns immediately (e.g., listener error), exit.
		if err != nil {
			log.Fatalf("gRPC serve error: %v", err)
		}
		return
	}

	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop() // stops accepting new conns/RPCs; waits for in-flight RPCs
		close(done)
	}()

	select {
	case <-done:
		log.Println("graceful shutdown complete")
	case <-ctx.Done():
		log.Printf("graceful shutdown timed out after %s; forcing stop", timeout)
		grpcServer.Stop()
	}

	_ = lis.Close()

}
