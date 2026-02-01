package grpc

import (
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip/v1"
	"ride-sharing/shared/types"
)

func protoToCoordinate(c *pb.Coordinate) *types.Coordinate {
	return &types.Coordinate{
		Latitude:  c.Latitude,
		Longitude: c.Longitude,
	}
}

func osrmToProtoRoute(o *types.OsrmApiResponse) *pb.Route {
	route := o.Routes[0]
	geometry := route.Geometry.Coordinates
	coordinates := make([]*pb.Coordinate, len(geometry))
	for i, coord := range geometry {
		coordinates[i] = &pb.Coordinate{
			Latitude:  coord[0],
			Longitude: coord[1],
		}
	}

	return &pb.Route{
		Geometry: []*pb.Geometry{
			{
				Coordinates: coordinates,
			},
		},
		Distance: route.Distance,
		Duration: route.Duration,
	}
}

func ToProtoRideFares(fares []*domain.RideFareModel) []*pb.RideFare {
	protoFares := make([]*pb.RideFare, len(fares))
	for i, fare := range fares {
		protoFares[i] = &pb.RideFare{
			Id:                fare.ID.Hex(),
			UserID:            fare.UserID,
			PackageSlug:       fare.PackageSlug,
			TotalPriceInCents: fare.TotalPriceInCents,
		}
	}
	return protoFares
}
