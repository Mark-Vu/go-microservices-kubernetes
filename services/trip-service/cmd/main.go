package main

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	service := service.NewTripService(inmemRepo)
	ctx := context.Background()

	fare := &domain.RideFareModel{
		UserID: "42",
	}
	t, err := service.CreateTrip(ctx, fare)
	if err != nil {
		log.Println(err)
	}
	log.Println(t)

	for {
		time.Sleep(1 * time.Second)
	}
}
