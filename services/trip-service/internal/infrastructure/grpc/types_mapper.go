package grpc

import (
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
