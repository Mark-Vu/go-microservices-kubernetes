package domain

import (
	"context"
)

type DriverRepository interface {
	CreateDriver(ctx context.Context, driver *DriverModel) (*DriverModel, error)
}

type DriverService interface {
	RegisterDriver(ctx context.Context, driverID, packageSlug string) (*DriverModel, error)
}
