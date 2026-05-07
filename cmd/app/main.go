package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/ElliAbby/go_cinema_system/internal/config"
	"github.com/ElliAbby/go_cinema_system/internal/cinemaService/usecase"
	"github.com/ElliAbby/go_cinema_system/internal/cinemaService/repository"
	cinemaHttp "github.com/ElliAbby/go_cinema_system/internal/cinemaService/transport/http"
	"github.com/ElliAbby/go_cinema_system/internal/db/postgres"

)


func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Ошибка при загрузке конфигурации: %v", err)
		return
	}
	log.Printf("Конфигурация сервера: %+v", cfg.Server)

	db, err := postgres.NewPostgresDB(&cfg.DB)
	if err != nil {
		log.Printf("Не удалось создать SQL-соединение: %v", err)
		return
	}
	defer db.Close()
	log.Println("Соединение с БД успешно установлено")

	repo := repository.New(db)
	uc := usecase.New(repo)
	handler := cinemaHttp.New(uc)
	routers := cinemaHttp.RegisterRouters(handler)

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
