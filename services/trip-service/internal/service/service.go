package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	tripv1 "ride-sharing/shared/proto/trip/v1"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripService struct {
	repo domain.TripRepository
}

func NewTripService(repo domain.TripRepository) *TripService {
	return &TripService{repo: repo}
}

func (s *TripService) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	// TODO: add driver selection logic
	newTrip := &domain.TripModel{
		ID:            primitive.NewObjectID(),
		UserID:        fare.UserID,
		Status:        "pending",
		RideFareModel: fare,
		Driver:        &tripv1.TripDriver{},
	}

	return s.repo.CreateTrip(ctx, newTrip)
}

func (s *TripService) GetRoute(ctx context.Context, pickup *types.Coordinate, destination *types.Coordinate) (*types.OsrmApiResponse, error) {
	url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM API :%v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response :%v", err)
	}

	var routeResp types.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &routeResp, nil
}

func (s *TripService) EstimatePackagesPriceWithRoute(ctx context.Context, route *types.OsrmApiResponse) ([]*domain.RideFareModel, error) {
	baseFares := getBaseFares()
	estimatedFares := make([]*domain.RideFareModel, len(baseFares))

	for i, fare := range baseFares {
		estimatedFares[i] = estimateFarePriceByRoute(fare, route)
	}

	return estimatedFares, nil
}

func (t *TripService) GenerateTripFares(ctx context.Context, fares []*domain.RideFareModel, userId string, route *types.OsrmApiResponse) ([]*domain.RideFareModel, error) {
	savedFares := make([]*domain.RideFareModel, 0, len(fares))

	for _, fare := range fares {
		id := primitive.NewObjectID()
		newFare := &domain.RideFareModel{
			ID:                id,
			UserID:            userId,
			PackageSlug:       fare.PackageSlug,
			TotalPriceInCents: fare.TotalPriceInCents,
			Route:             route,
		}
		if err := t.repo.SaveRideFare(ctx, newFare); err != nil {
			return nil, fmt.Errorf("failed to save ride fare: %v", err)
		}
		savedFares = append(savedFares, newFare)
	}

	return savedFares, nil
}

func (t *TripService) GetRideFareByID(ctx context.Context, fareId string, userId string) (*domain.RideFareModel, error) {
	fareIdObj, err := primitive.ObjectIDFromHex(fareId)
	if err != nil {
		return nil, fmt.Errorf("failed to convert fare id to object id: %v", err)
	}
	fare, err := t.repo.GetRideFareByID(ctx, fareIdObj)
	if err != nil {
		return nil, fmt.Errorf("failed to get ride fare: %v", err)
	}
	if fare.UserID != userId {
		return nil, fmt.Errorf("ride fare not found")
	}
	return fare, nil
}

func estimateFarePriceByRoute(fare *domain.RideFareModel, route *types.OsrmApiResponse) *domain.RideFareModel {
	pricing := DefaultPricingConfig()
	distanceKm := route.Routes[0].Distance
	durationMinutes := route.Routes[0].Duration
	vehicleSpecificPrice := fare.TotalPriceInCents

	distanceFare := distanceKm * pricing.PricePerUnitOfDistance
	timeFare := durationMinutes * pricing.PricingPerMinute

	totalFare := distanceFare + timeFare + vehicleSpecificPrice

	return &domain.RideFareModel{
		PackageSlug:       fare.PackageSlug,
		TotalPriceInCents: totalFare,
	}
}
func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 200,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 350,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 400,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000,
		},
	}
}
