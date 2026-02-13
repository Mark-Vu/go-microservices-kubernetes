package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/services/driver-service/internal/service"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	"github.com/rabbitmq/amqp091-go"
)

type TripConsumer struct {
	rabbitmq      *messaging.RabbitMQ
	driverService *service.DriverService
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ, driverService *service.DriverService) *TripConsumer {
	return &TripConsumer{
		rabbitmq:      rabbitmq,
		driverService: driverService,
	}
}

func (c *TripConsumer) Start(ctx context.Context) error {
	log.Println("Starting trip.event.created consumer...")
	return c.rabbitmq.Consume(ctx, messaging.FindAvailableDriversQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		return c.handleTripCreated(ctx, msg)
	})
}

func (c *TripConsumer) handleTripCreated(ctx context.Context, msg amqp091.Delivery) error {
	log.Printf("Received trip created event: %s", string(msg.Body))
	var tripEvent contracts.AmqpMessage
	if err := json.Unmarshal(msg.Body, &tripEvent); err != nil {
		log.Printf("Failed to unmarshal trip event: %v", err)
		return fmt.Errorf("failed to unmarshal trip event: %w", err)
	}
	var payload messaging.TripEventData
	if err := json.Unmarshal(tripEvent.Data, &payload); err != nil {
		log.Printf("Failed to unmarshal trip event data: %v", err)
		return fmt.Errorf("failed to unmarshal trip event data: %w", err)
	}

	switch msg.RoutingKey {
	case contracts.TripEventCreated:
		log.Printf("Trip created event received: %+v", payload)

		return c.handleFindAndNotifyDrivers(ctx, payload)
	}
	return nil
}

func (c *TripConsumer) Close() error {
	return c.rabbitmq.Close()
}

func (c *TripConsumer) handleFindAndNotifyDrivers(ctx context.Context, payload messaging.TripEventData) error {
	matchedDriverIDs := c.driverService.FindAvailableDrivers(ctx, payload.Trip.SelectedFare.PackageSlug)

	if len(matchedDriverIDs) == 0 {
		message := contracts.AmqpMessage{
			OwnerID: payload.Trip.UserID,
			Data:    nil,
		}
		// Notify riders that no available drivers found
		if err := c.rabbitmq.Publish(ctx, contracts.TripEventNoDriversFound, message); err != nil {
			return fmt.Errorf("failed to publish no drivers found event: %v", err)
		}
		log.Printf("Published no drivers found event to %s with body %s", messaging.TripExchange, message)
		return nil
	}
	foundDriverID := matchedDriverIDs[0]
	log.Printf("Found driver ID: %s", foundDriverID)
	marshalledData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	message := contracts.AmqpMessage{
		OwnerID: foundDriverID,
		Data:    marshalledData,
	}
	// Notify the driver about a potential trip
	if err := c.rabbitmq.Publish(ctx, contracts.DriverCmdTripRequest, message); err != nil {
		log.Printf("Failed to publish message to exchange: %v", err)
		return err
	}
	log.Printf("Published trip request event to %s with body %s", contracts.DriverCmdTripRequest, message)

	return nil
}
