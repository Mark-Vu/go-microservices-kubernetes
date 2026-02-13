package events

import (
	"context"
	"log"
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

	// TODO: Unmarshal and process
	// var event TripCreatedEvent
	// if err := json.Unmarshal(msg.Body, &event); err != nil {
	//     return err
	// }

	// TODO: Business logic (assign driver, etc.)

	return nil
}

func (c *TripConsumer) Close() error {
	return c.rabbitmq.Close()
}
