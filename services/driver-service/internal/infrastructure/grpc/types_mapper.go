package grpc

import (
	"ride-sharing/services/driver-service/internal/domain"
	pb "ride-sharing/shared/proto/driver/v1"
)

func ToProtoDriver(driver *domain.DriverModel) *pb.Driver {
	return &pb.Driver{
		Id:             driver.ID.Hex(),
		Name:           driver.Name,
		ProfilePicture: driver.ProfilePicture,
		CarPlate:       driver.CarPlate,
		Geohash:        driver.Geohash,
		PackageSlug:    driver.PackageSlug,
		Location:       driver.Location,
	}
}
