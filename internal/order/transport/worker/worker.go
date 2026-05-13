package worker

import (
	"context"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/ElliAbby/go_cinema_system/internal/order"

)

type WorkerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type Worker struct {
	reader *kafka.Reader
	uc     order.UseCase
}

func New(cfg WorkerConfig, uc order.UseCase) *Worker {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Brokers,
		Topic:       cfg.Topic,
		GroupID:     cfg.GroupID,
		MinBytes:    1,
		MaxBytes:    10e6,
		StartOffset: kafka.LastOffset,
	})

	return &Worker{
		reader: reader,
		uc:     uc,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	log.Println("Order worker started, listening for payment messages")

	for {
		msg, err := w.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("Order worker stopped")
				return nil
			}
			log.Printf("Kafka fetch error: %v", err)
			continue
		}

		if err := handlePaymentMessage(ctx, w.uc, msg.Value); err != nil {
			log.Printf("Payment message processing error: %v", err)
			continue
		}

		if err := w.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Kafka commit error: %v", err)
		}
	}
}

func (w *Worker) Close() error {
	if w.reader != nil {
		return w.reader.Close()
	}
	return nil
}
