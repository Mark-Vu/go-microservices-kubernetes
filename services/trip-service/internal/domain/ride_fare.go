package domain

import (
	pb "ride-sharing/shared/proto/trip/v1"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                primitive.ObjectID
	UserID            string
	PackageSlug       string // ex: van, luxury, sedan
	TotalPriceInCents float64
	Route             *types.OsrmApiResponse
}

func (f *RideFareModel) ToProto() *pb.RideFare {
	if f == nil {
		return nil
	}
	return &pb.RideFare{
		Id:                f.ID.Hex(),
		UserID:            f.UserID,
		PackageSlug:       f.PackageSlug,
		TotalPriceInCents: f.TotalPriceInCents,
	}
}
