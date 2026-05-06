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
)


func main() {
	cfg := config.Load()
	log.Printf("Конфигурация сервера: %+v", cfg.Server)

	repo := repository.New()
	uc := usecase.New(repo)
	handler := cinemaHttp.New(uc)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("Добро пожаловать в наш кинотеатр!"))
	})
	mux.HandleFunc("/test", handler.TestEndpoint)
	mux.HandleFunc("/slow", handler.SlowEndpoint)

	srv := &http.Server{
		Addr: cfg.Server.Addr,
		Handler: mux,
		ReadTimeout: cfg.Server.ReadTimeout,
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
	// для двойного Ctrl+C
	stop()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutDownCtx); err != nil {
		log.Fatalf("Ошибка при принудительном выключении сервера: %v", err)
	}
	log.Println("Сервер успешно остановлен")
}
