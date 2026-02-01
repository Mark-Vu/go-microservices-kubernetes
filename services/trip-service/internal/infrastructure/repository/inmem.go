package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *inmemRepository) GetRideFareByID(ctx context.Context, id primitive.ObjectID) (*domain.RideFareModel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	fare, ok := r.rideFares[id.Hex()]
	if !ok {
		return nil, fmt.Errorf("ride fare not found")
	}
	return fare, nil
}
