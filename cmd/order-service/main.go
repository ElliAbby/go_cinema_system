package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	orderRepo "github.com/ElliAbby/go_cinema_system/internal/order/repository"
	orderUseCase "github.com/ElliAbby/go_cinema_system/internal/order/usecase"
	"github.com/ElliAbby/go_cinema_system/internal/order/transport/worker"
	"github.com/ElliAbby/go_cinema_system/internal/platform/config"
	"github.com/ElliAbby/go_cinema_system/internal/platform/db/postgres"
	"github.com/ElliAbby/go_cinema_system/internal/platform/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	metrics.InitMetrics(db, "order-service")
	log.Println("Order service business metrics initialized")
	metrics.InitKafkaMetrics("order-service", cfg.Kafka.BookingPaymentsTopic)
	log.Println("Order service metrics initialized")

	metricsServer := &http.Server{
		Addr:    cfg.MetricsAddr,
		Handler: promhttp.Handler(),
	}

	go func() {
		log.Printf("Order service metrics exposed on %s", cfg.MetricsAddr)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	w := worker.New(worker.WorkerConfig{
		Brokers:     cfg.Kafka.Brokers,
		Topic:       cfg.Kafka.BookingPaymentsTopic,
		GroupID:     cfg.Kafka.BookingPaymentsGroup,
		ServiceName: "order-service",
	}, uc)
	defer w.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := w.Start(ctx); err != nil {
		log.Fatalf("Worker error: %v", err)
	}

	shutDownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := metricsServer.Shutdown(shutDownCtx); err != nil {
		log.Printf("Metrics server shutdown error: %v", err)
	}
}