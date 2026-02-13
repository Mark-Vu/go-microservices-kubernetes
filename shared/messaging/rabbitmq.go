package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/shared/contracts"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

const (
	TripExchange    string = "trip-exchange"
	DriverExchange  string = "driver-exchange"
	PaymentExchange string = "payment-exchange"
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
	err := r.ch.ExchangeDeclare(
		TripExchange, // exchange name
		"topic",      // exchange type
		true,         // durable
		false,        // auto-delete
		false,        // no-wait
		false,        // internal
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	if err := r.declareAndBindQueue(
		FindAvailableDriversQueue, // Create a queue for finding available drivers
		[]string{ // queue will listen to these routing keys
			contracts.TripEventCreated,
			contracts.TripEventDriverNotInterested,
		},
		TripExchange, // bind to this exchange
	); err != nil {
		return fmt.Errorf("failed to declare and bind queue: %w", err)
	}
	return nil
}

func (r *RabbitMQ) declareAndBindQueue(queueName string, routingKeys []string, exchange string) error {
	queue, err := r.ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	for _, routingKey := range routingKeys {
		if err := r.ch.QueueBind(
			queue.Name, // queue name
			routingKey, // routing key
			exchange,   // exchange name
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to bind queue: %w", err)
		}
	}
	return nil
}

func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, body contracts.AmqpMessage) error {
	log.Printf("Publishing message to %s with routing key %s", TripExchange, routingKey)
	message, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}
	return r.ch.PublishWithContext(ctx,
		TripExchange, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp091.Persistent,
		},
	)
}

type MessageHandler func(ctx context.Context, msg amqp091.Delivery) error

func (r *RabbitMQ) Consume(ctx context.Context, queueName string, handlerFunc MessageHandler) error {
	if err := r.ch.Qos(
		1,     // prefetch count - max unacked messages per consumer
		0,     // prefetch size - 0 means no limit by bytes
		false, // global - Apply prefetchCount to all consumers on this channel
	); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := r.ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}
			if err := handlerFunc(ctx, msg); err != nil {
				log.Printf("failed to process message: %v", err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}
}

func (r *RabbitMQ) Close() error {
	if r.conn == nil {
		return nil // idempotent close
	}
	return r.conn.Close()
}
