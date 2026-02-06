package domain

import (
	pb "ride-sharing/shared/proto/driver/v1"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DriverModel struct {
	ID          primitive.ObjectID
	DriverID    string
	PackageSlug string
}

type driverInMap struct {
	Driver *pb.Driver
}
