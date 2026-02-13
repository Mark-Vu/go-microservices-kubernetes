package service

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"ride-sharing/services/driver-service/internal/domain"
	"ride-sharing/services/driver-service/internal/fixtures"
	pb "ride-sharing/shared/proto/driver/v1"
	"ride-sharing/shared/util"
	"sync"

	"github.com/mmcloughlin/geohash"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DriverStatus represents the current state of a driver
type DriverStatus string

const (
	DriverStatusAvailable DriverStatus = "available" // Ready to accept trips
	DriverStatusBusy      DriverStatus = "busy"      // Currently on a trip
)

// OnlineDriver represents a driver who is currently online
type OnlineDriver struct {
	Driver *pb.Driver
	Status DriverStatus
}

// DriverService manages driver operations
type DriverService struct {
	repo          domain.DriverRepository
	onlineDrivers map[string]*OnlineDriver // driverID -> OnlineDriver
	mu            sync.RWMutex             // Protects onlineDrivers map
}

// NewDriverService creates a new driver service instance
func NewDriverService(repo domain.DriverRepository) *DriverService {
	return &DriverService{
		repo:          repo,
		onlineDrivers: make(map[string]*OnlineDriver),
	}
}

// RegisterDriver adds a driver to the online pool (driver goes "online")
func (s *DriverService) RegisterDriver(ctx context.Context, driverID, packageSlug string) (*pb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if already online
	if _, exists := s.onlineDrivers[driverID]; exists {
		return nil, fmt.Errorf("driver %s is already online", driverID)
	}

	// Pick random route and location
	randomIndex := rand.IntN(len(fixtures.PredefinedRoutes))
	randomRoute := fixtures.PredefinedRoutes[randomIndex]
	startLocation := randomRoute[0]

	// Generate geohash for proximity search
	driverGeohash := geohash.Encode(startLocation[0], startLocation[1])

	// Create driver
	driver := &pb.Driver{
		Id:             driverID,
		Geohash:        driverGeohash,
		Location:       &pb.Location{Latitude: startLocation[0], Longitude: startLocation[1]},
		Name:           "Lando Norris", // TODO: Get from user profile
		PackageSlug:    packageSlug,
		ProfilePicture: util.GetRandomAvatar(randomIndex),
		CarPlate:       fixtures.GenerateRandomPlate(),
	}

	// Add to online pool
	s.onlineDrivers[driverID] = &OnlineDriver{
		Driver: driver,
		Status: DriverStatusAvailable,
	}

	log.Printf("Driver %s registered online with package %s", driverID, packageSlug)
	return driver, nil
}

// UnregisterDriver removes a driver from the online pool (driver goes "offline")
func (s *DriverService) UnregisterDriver(ctx context.Context, driverID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.onlineDrivers[driverID]; !exists {
		return fmt.Errorf("driver %s is not online", driverID)
	}

	delete(s.onlineDrivers, driverID)
	log.Printf("Driver %s unregistered (offline)", driverID)
	return nil
}

// FindAvailableDrivers returns IDs of available drivers matching package type
func (s *DriverService) FindAvailableDrivers(ctx context.Context, packageSlug string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var availableDrivers []string

	for driverID, onlineDriver := range s.onlineDrivers {
		// Only match available drivers with correct package
		if onlineDriver.Status == DriverStatusAvailable &&
			onlineDriver.Driver.PackageSlug == packageSlug {
			availableDrivers = append(availableDrivers, driverID)
		}
	}

	log.Printf("Found %d available drivers for package %s", len(availableDrivers), packageSlug)
	return availableDrivers
}

// SetDriverBusy marks a driver as busy (on a trip)
func (s *DriverService) SetDriverBusy(ctx context.Context, driverID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	driver, exists := s.onlineDrivers[driverID]
	if !exists {
		return fmt.Errorf("driver %s is not online", driverID)
	}

	driver.Status = DriverStatusBusy
	log.Printf("Driver %s set to busy", driverID)
	return nil
}

// SetDriverAvailable marks a driver as available (trip completed)
func (s *DriverService) SetDriverAvailable(ctx context.Context, driverID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	driver, exists := s.onlineDrivers[driverID]
	if !exists {
		return fmt.Errorf("driver %s is not online", driverID)
	}

	driver.Status = DriverStatusAvailable
	log.Printf("Driver %s set to available", driverID)
	return nil
}

// GetOnlineDriver returns a driver if they're online
func (s *DriverService) GetOnlineDriver(ctx context.Context, driverID string) (*pb.Driver, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	driver, exists := s.onlineDrivers[driverID]
	if !exists {
		return nil, fmt.Errorf("driver %s is not online", driverID)
	}

	return driver.Driver, nil
}

// GetOnlineDriverCount returns the number of online drivers
func (s *DriverService) GetOnlineDriverCount(ctx context.Context) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.onlineDrivers)
}

// --- Database operations (for persistence) ---

// CreateDriverProfile creates a driver profile in the database
func (s *DriverService) CreateDriverProfile(ctx context.Context, driverID, packageSlug string) (*domain.DriverModel, error) {
	// This saves to database (permanent storage)
	driver := &domain.DriverModel{
		ID:             primitive.NewObjectID(),
		DriverID:       driverID,
		PackageSlug:    packageSlug,
		Name:           "John Doe",
		ProfilePicture: util.GetRandomAvatar(1),
		CarPlate:       fixtures.GenerateRandomPlate(),
	}

	createdDriver, err := s.repo.CreateDriver(ctx, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver profile: %w", err)
	}

	log.Printf("Driver profile created in database: %s", driverID)
	return createdDriver, nil
}
