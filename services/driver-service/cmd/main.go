package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/shared/env"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8083")
)

func main() {
	lis, err := net.Listen("tcp", httpAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Starting grpc server
	grpcServer := grpc.NewServer()

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
