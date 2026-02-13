package domain

import (
	"context"
	pb "ride-sharing/shared/proto/driver/v1"
)

// DriverRepository handles database operations for driver profiles
type DriverRepository interface {
	CreateDriver(ctx context.Context, driver *DriverModel) (*DriverModel, error)
}

// DriverService handles business logic for driver operations
type DriverService interface {
	// Online/Offline management
	RegisterDriver(ctx context.Context, driverID, packageSlug string) (*pb.Driver, error)
	UnregisterDriver(ctx context.Context, driverID string) error

	// Trip matching
	FindAvailableDrivers(ctx context.Context, packageSlug string) []string

	// Status management
	SetDriverBusy(ctx context.Context, driverID string) error
	SetDriverAvailable(ctx context.Context, driverID string) error

	// Query
	GetOnlineDriver(ctx context.Context, driverID string) (*pb.Driver, error)
	GetOnlineDriverCount(ctx context.Context) int

	// Database operations
	CreateDriverProfile(ctx context.Context, driverID, packageSlug string) (*DriverModel, error)
}
