package repository

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
	"sync"
)

type inmemRepository struct {
	mu        sync.Mutex
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *inmemRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func (r *inmemRepository) SaveRideFare(ctx context.Context, fare *domain.RideFareModel) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.rideFares[fare.ID.Hex()] = fare
	return nil
}
