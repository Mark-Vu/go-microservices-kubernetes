package messaging

import pb "ride-sharing/shared/proto/trip/v1"

const (
	FindAvailableDriversQueue = "find_available_drivers"
)

type TripEventData struct {
	Trip *pb.Trip `json:"trip"`
}
