package messaging

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp091.Connection
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to dial RabbitMQ: %w", err)
	}
	return &RabbitMQ{conn: conn}, nil
}

func (r *RabbitMQ) Close() error {
	if r.conn == nil {
		return fmt.Errorf("connection is not open")
	}
	r.conn.Close()
	return nil
}
