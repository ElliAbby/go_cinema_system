package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	orderRepo "github.com/ElliAbby/go_cinema_system/internal/order/repository"
	orderUseCase "github.com/ElliAbby/go_cinema_system/internal/order/usecase"
	"github.com/ElliAbby/go_cinema_system/internal/platform/config"
	"github.com/ElliAbby/go_cinema_system/internal/platform/db/postgres"
	"github.com/ElliAbby/go_cinema_system/internal/order/transport/worker"

)

func main() {
	cfg, err := config.LoadWorkerConfig()
	if err != nil {
		log.Fatalf("Failed to load worker config: %v", err)
	}

	db, err := postgres.NewPostgresDB(&cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Order Service: connected to database")

	repo := orderRepo.New(db)
	uc := orderUseCase.New(repo)

	w := worker.New(worker.WorkerConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.BookingPaymentsTopic,
		GroupID: cfg.Kafka.BookingPaymentsGroup,
	}, uc)
	defer w.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := w.Start(ctx); err != nil {
		log.Fatalf("Worker error: %v", err)
	}
}