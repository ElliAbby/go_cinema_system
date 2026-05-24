package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/ElliAbby/go_cinema_system/internal/platform/config"
	"github.com/ElliAbby/go_cinema_system/internal/cinema/usecase"
	"github.com/ElliAbby/go_cinema_system/internal/cinema/repository"
	cinemaHttp "github.com/ElliAbby/go_cinema_system/internal/cinema/transport/http"
	"github.com/ElliAbby/go_cinema_system/internal/platform/db/postgres"
	cinemaKafka "github.com/ElliAbby/go_cinema_system/internal/platform/kafka"
	authjwt "github.com/ElliAbby/go_cinema_system/internal/platform/auth/jwt"
	"github.com/ElliAbby/go_cinema_system/internal/platform/metrics"

)


func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Ошибка при загрузке конфигурации: %v", err)
		return
	}
	log.Printf("Конфигурация сервера: %+v", cfg.Server)

	tokenManager := authjwt.New(cfg.JWT.SecretKey)

	db, err := postgres.NewPostgresDB(&cfg.DB)
	if err != nil {
		log.Printf("Не удалось создать SQL-соединение: %v", err)
		return
	}
	defer db.Close()
	log.Println("Соединение с БД успешно установлено")

	publisher, err := cinemaKafka.NewPublisher(cfg.Kafka.Brokers, cfg.Kafka.BookingPaymentsTopic)
	if err != nil {
		log.Printf("Не удалось создать Kafka publisher: %v", err)
		return
	}
	defer publisher.Close()
	log.Println("Kafka publisher инициализирован")

	if err := postgres.SeedDatabase(db); err != nil {
		log.Printf("Ошибка при загрузке seed-данных: %v", err)
		return
	}
	log.Println("Seed-данные успешно загружены (если требовалось)")

	metrics.InitMetrics(db, "cinema-service")
	metrics.InitKafkaMetrics("cinema-service", cfg.Kafka.BookingPaymentsTopic)
	log.Println("Метрики инициализированы")

	repo := repository.New(db, cfg.ReservationHoldDuration)
	uc := usecase.New(repo, publisher, tokenManager, cfg.ReservationHoldDuration)
	handler := cinemaHttp.New(uc)
	routers := cinemaHttp.RegisterRouters(handler, tokenManager)

	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      routers,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("Сервер запущен на порту %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("Получен сигнал завершения, выключаем сервер...")

	shutDownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutDownCtx); err != nil {
		log.Fatalf("Ошибка при принудительном выключении сервера: %v", err)
	}
	log.Println("Сервер успешно остановлен")
}
