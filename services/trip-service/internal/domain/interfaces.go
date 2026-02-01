package domain

import (
	"context"
	"ride-sharing/shared/types"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, fare *RideFareModel) error
	GetRideFareByID(ctx context.Context, id string) (*RideFareModel, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup *types.Coordinate, destination *types.Coordinate) (*types.OsrmApiResponse, error)
	EstimatePackagesPriceWithRoute(ctx context.Context, route *types.OsrmApiResponse) ([]*RideFareModel, error)
	GenerateTripFares(ctx context.Context, fares []*RideFareModel, userId string) ([]*RideFareModel, error)
	GetRideFareByID(ctx context.Context, fareId string, userId string) (*RideFareModel, error)
}
