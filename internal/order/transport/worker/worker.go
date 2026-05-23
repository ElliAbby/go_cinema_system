package worker

import (
	"context"
	"errors"
	"log"

	"github.com/ElliAbby/go_cinema_system/internal/platform/metrics"
	"github.com/segmentio/kafka-go"

	"github.com/ElliAbby/go_cinema_system/internal/order"

)

type WorkerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
	ServiceName string
}

type Worker struct {
	reader      *kafka.Reader
	uc          order.UseCase
	topic       string
	serviceName string
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
		reader:      reader,
		uc:          uc,
		topic:       cfg.Topic,
		serviceName: cfg.ServiceName,
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
			metrics.IncKafkaConsumerErrors(w.topic, "fetch_error")
			continue
		}

		metrics.IncKafkaMessagesReceived(w.topic)

		if err := handlePaymentMessage(ctx, w.topic, w.serviceName, w.uc, msg.Value); err != nil {
			log.Printf("Payment message processing error: %v", err)
			metrics.IncKafkaMessagesFailed(w.topic, "processing_error")
			continue
		}

		metrics.IncKafkaMessagesProcessed(w.topic)

		if err := w.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Kafka commit error: %v", err)
			metrics.IncKafkaConsumerErrors(w.topic, "commit_error")
		}
	}
}

func (w *Worker) Close() error {
	if w.reader != nil {
		return w.reader.Close()
	}
	return nil
}
