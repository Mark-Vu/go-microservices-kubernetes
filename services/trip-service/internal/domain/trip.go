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

func (t *TripModel) ToProto() *pb.Trip {
	if t == nil {
		return nil
	}

	var protoRoute *pb.Route
	if t.RideFareModel != nil && t.RideFareModel.Route != nil {
		protoRoute = t.RideFareModel.Route.ToProto()
	}

	return &pb.Trip{
		Id:           t.ID.Hex(),
		UserID:       t.UserID,
		Status:       t.Status,
		SelectedFare: t.RideFareModel.ToProto(),
		Driver:       t.Driver,
		Route:        protoRoute,
	}
}
