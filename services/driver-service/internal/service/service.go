package service

import (
	"context"
	"fmt"
	"ride-sharing/services/driver-service/internal/domain"
	"ride-sharing/services/driver-service/internal/fixtures"
	pb "ride-sharing/shared/proto/driver/v1"
	"ride-sharing/shared/util"

	"math/rand/v2"

	"github.com/mmcloughlin/geohash"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DriverService struct {
	repo domain.DriverRepository
}

func NewDriverService(repo domain.DriverRepository) *DriverService {
	return &DriverService{repo: repo}
}

func (s *DriverService) RegisterDriver(ctx context.Context, driverID, packageSlug string) (*domain.DriverModel, error) {
	randomIndex := rand.IntN(len(fixtures.PredefinedRoutes))
	randomRoute := fixtures.PredefinedRoutes[randomIndex]

	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	driver := &domain.DriverModel{
		ID:          primitive.NewObjectID(),
		DriverID:    driverID,
		PackageSlug: packageSlug,
		Geohash:     geohash,
		Location: &pb.Location{
			Latitude:  randomRoute[0][0],
			Longitude: randomRoute[0][1],
		},
		Name:           "John Doe",
		ProfilePicture: util.GetRandomAvatar(1),
		CarPlate:       fixtures.GenerateRandomPlate(),
	}
	createdDriver, err := s.repo.CreateDriver(ctx, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %v", err)
	}
	return createdDriver, nil
}
