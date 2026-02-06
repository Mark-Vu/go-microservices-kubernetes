package repository

import (
	"context"
	"ride-sharing/services/driver-service/internal/domain"
	"sync"
)

type inmemRepository struct {
	mu      sync.Mutex
	drivers map[string]*domain.DriverModel
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		drivers: make(map[string]*domain.DriverModel),
	}
}

func (r *inmemRepository) CreateDriver(ctx context.Context, driver *domain.DriverModel) (*domain.DriverModel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.drivers[driver.ID.Hex()] = driver
	return driver, nil
}
