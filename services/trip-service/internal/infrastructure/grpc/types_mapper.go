package grpc

import (
	pb "ride-sharing/shared/proto/trip/v1"
	"ride-sharing/shared/types"
)

// protoToCoordinate converts protobuf Coordinate to domain Coordinate
func protoToCoordinate(c *pb.Coordinate) *types.Coordinate {
	if c == nil {
		return nil
	}
	return &types.Coordinate{
		Latitude:  c.Latitude,
		Longitude: c.Longitude,
	}
}
