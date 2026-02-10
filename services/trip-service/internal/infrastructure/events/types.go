package events

type TripCreatedEvent struct {
	TripID      string  `json:"tripId"`
	UserID      string  `json:"userId"`
	Status      string  `json:"status"`
	PackageSlug string  `json:"packageSlug"`
	FarePrice   float64 `json:"farePrice"`
}
