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
		log.Fatalf("Не удалось загрузить конфигурацию воркера: %v", err)
	}

	db, err := postgres.NewPostgresDB(&cfg.DB)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()
	log.Println("Сервис заказов: подключение к базе данных установлено")

	repo := orderRepo.New(db)
	uc := orderUseCase.New(repo)

	metrics.InitMetrics(db, "order-service")
	log.Println("Бизнес-метрики сервиса заказов инициализированы")
	metrics.InitKafkaMetrics("order-service", cfg.Kafka.BookingPaymentsTopic)
	log.Println("Метрики сервиса заказов инициализированы")

	metricsServer := &http.Server{
		Addr:    cfg.MetricsAddr,
		Handler: promhttp.Handler(),
	}

	go func() {
		log.Printf("Метрики сервиса заказов доступны на %s", cfg.MetricsAddr)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка сервера метрик: %v", err)
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
		log.Fatalf("Ошибка воркера: %v", err)
	}

	shutDownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := metricsServer.Shutdown(shutDownCtx); err != nil {
		log.Printf("Ошибка при остановке сервера метрик: %v", err)
	}
}