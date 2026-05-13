package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type PaymentRequestedEvent struct {
	BookingID   string    `json:"booking_id"`
	UserID      int       `json:"user_id"`
	TotalPrice  float64   `json:"total_price"`
	RequestedAt time.Time `json:"requested_at"`
}

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(brokers []string, topic string) (*Publisher, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("kafka brokers are required")
	}
	if topic == "" {
		return nil, fmt.Errorf("kafka topic is required")
	}

	return &Publisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			RequiredAcks: kafka.RequireOne,
			Async:        false,
		},
	}, nil
}

func (p *Publisher) Close() error {
	if p == nil || p.writer == nil {
		return nil
	}
	return p.writer.Close()
}

func (p *Publisher) PublishPaymentRequested(ctx context.Context, event PaymentRequestedEvent) error {
	if p == nil || p.writer == nil {
		return fmt.Errorf("kafka publisher is not initialized")
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.BookingID),
		Value: body,
	})
}