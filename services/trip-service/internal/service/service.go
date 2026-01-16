package service

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripService struct {
	repo domain.TripRepository
}

func NewTripService(repo domain.TripRepository) *TripService {
	return &TripService{repo: repo}
}

func (s *TripService) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	newTrip := &domain.TripModel{
		ID:            primitive.NewObjectID(),
		UserID:        "",
		Status:        "",
		RideFareModel: fare,
	}

	return s.repo.CreateTrip(ctx, newTrip)
}
