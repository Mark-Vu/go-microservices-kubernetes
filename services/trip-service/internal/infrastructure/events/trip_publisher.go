package events

import (
	"context"
	"encoding/json"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/messaging"
)

type TripEventPublisher struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripPublisher(rabbitmq *messaging.RabbitMQ) *TripEventPublisher {
	return &TripEventPublisher{rabbitmq: rabbitmq}
}

func (p *TripEventPublisher) PublishTripCreated(ctx context.Context, trip *domain.TripModel) error {
	event := TripCreatedEvent{
		TripID:      trip.ID.Hex(),
		UserID:      trip.UserID,
		Status:      trip.Status,
		PackageSlug: trip.RideFareModel.PackageSlug,
		FarePrice:   trip.RideFareModel.TotalPriceInCents,
	}
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return p.rabbitmq.Publish(ctx, "trip.event.created", body)
}

func (p *TripEventPublisher) Close() error {
	return p.rabbitmq.Close()
}
