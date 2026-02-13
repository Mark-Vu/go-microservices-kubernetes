package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	"github.com/rabbitmq/amqp091-go"
)

type TripConsumer struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ) *TripConsumer {
	return &TripConsumer{rabbitmq: rabbitmq}
}

func (c *TripConsumer) Start(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func() {
		log.Println("Starting trip.event.created consumer...")
		errChan <- c.rabbitmq.Consume(ctx, messaging.FindAvailableDriversQueue, c.handleTripCreated)
	}()

	return <-errChan
}

func (c *TripConsumer) handleTripCreated(msg amqp091.Delivery) error {
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

	log.Printf("Trip created event received: %+v", payload)
	return nil
}

func (c *TripConsumer) Close() error {
	return c.rabbitmq.Close()
}
