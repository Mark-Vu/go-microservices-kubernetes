package messaging

import (
	"context"
	"fmt"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to dial RabbitMQ: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}
	r := &RabbitMQ{
		conn: conn,
		ch:   ch,
	}

	if err := r.setupExchangesAndQueues("trip.event.created"); err != nil {
		_ = r.Close()
		return nil, fmt.Errorf("failed to setup exchanges and queues: %w", err)
	}

	return r, nil
}

func (r *RabbitMQ) setupExchangesAndQueues(name string) error {
	_, err := r.ch.QueueDeclare(
		name,  // queue name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	return nil
}

func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, body []byte) error {
	return r.ch.PublishWithContext(ctx,
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)
}

func (r *RabbitMQ) Close() error {
	if r.conn == nil {
		return nil // idempotent close
	}
	return r.conn.Close()
}
