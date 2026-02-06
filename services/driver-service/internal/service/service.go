package service

import (
	"context"
	"ride-sharing/services/driver-service/internal/domain"
)

type DriverService struct {
	repo domain.DriverRepository
}

func NewDriverService(repo domain.DriverRepository) *DriverService {
	return &DriverService{repo: repo}
}

func (s *DriverService) RegisterDriver(ctx context.Context, driver *domain.DriverModel) (*domain.DriverModel, error) {
	return s.repo.CreateDriver(ctx, driver)
}
