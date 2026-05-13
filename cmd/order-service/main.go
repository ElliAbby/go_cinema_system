package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/ElliAbby/go_cinema_system/internal/orderService"
	orderRepo "github.com/ElliAbby/go_cinema_system/internal/orderService/repository"
	orderUseCase "github.com/ElliAbby/go_cinema_system/internal/orderService/usecase"
	"github.com/ElliAbby/go_cinema_system/internal/config"
	"github.com/ElliAbby/go_cinema_system/internal/db/postgres"
	cinemaKafka "github.com/ElliAbby/go_cinema_system/internal/kafka"
	"github.com/segmentio/kafka-go"
)

func main() {
	cfg, err := config.LoadWorkerConfig()
	if err != nil {
		log.Printf("Не удалось загрузить конфигурацию worker: %v", err)
		return
	}

	db, err := postgres.NewPostgresDB(&cfg.DB)
	if err != nil {
		log.Printf("Не удалось создать SQL-соединение: %v", err)
		return
	}
	defer db.Close()
	log.Println("Order Service: connected to database")

	repo := orderRepo.New(db)
	uc := orderUseCase.New(repo)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Kafka.Brokers,
		Topic:       cfg.Kafka.BookingPaymentsTopic,
		GroupID:     cfg.Kafka.BookingPaymentsGroup,
		MinBytes:    1,
		MaxBytes:    10e6,
		StartOffset: kafka.LastOffset,
	})
	defer reader.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("Order service started")
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("Order service stopped")
				return
			}
			log.Printf("Kafka fetch error: %v", err)
			continue
		}

		if err := handlePaymentMessage(ctx, uc, msg.Value); err != nil {
			log.Printf("Payment message processing error: %v", err)
			continue
		}

		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Kafka commit error: %v", err)
		}
	}
}

func handlePaymentMessage(ctx context.Context, uc orderService.UseCase, payload []byte) error {
	var event cinemaKafka.PaymentRequestedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	start := time.Now()
	booking, tickets, err := uc.ProcessPayment(ctx, event.UserID, event.BookingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Booking %s already processed or not found", event.BookingID)
			return nil
		}

		log.Printf("Payment processing failed for booking %s: %v", event.BookingID, err)
		return nil
	}

	log.Printf("Booking %s paid successfully in %s, tickets=%d, total=%.2f", booking.ID, time.Since(start).String(), len(tickets), booking.TotalPrice)
	return nil
}