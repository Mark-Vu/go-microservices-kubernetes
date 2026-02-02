package domain

import (
	pb "ride-sharing/shared/proto/trip/v1"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID            primitive.ObjectID
	UserID        string
	Status        string
	RideFareModel *RideFareModel
	Driver        *pb.TripDriver
}
