package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
)

type TripEventPublisher struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripPublisher(rabbitmq *messaging.RabbitMQ) *TripEventPublisher {
	return &TripEventPublisher{rabbitmq: rabbitmq}
}

func (p *TripEventPublisher) PublishTripCreated(ctx context.Context, trip *domain.TripModel) error {
	payload := messaging.TripEventData{
		Trip: trip.ToProto(),
	}
	tripJson, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	body := contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    tripJson,
	}

	log.Printf("Publishing trip created event to %s with body %s", messaging.TripExchange, body)
	return p.rabbitmq.Publish(ctx, contracts.TripEventCreated, body)
}

func (p *TripEventPublisher) Close() error {
	return p.rabbitmq.Close()
}
