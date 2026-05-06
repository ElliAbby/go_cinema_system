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

	mux := http.NewServeMux()

	// Корневой путь
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Добро пожаловать в кинотеатр API!"}`))
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*1000*1000*1000) // 2 seconds
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"unhealthy","error":"database unavailable"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Legacy test endpoints
	mux.HandleFunc("/test", handler.TestEndpoint)
	mux.HandleFunc("/slow", handler.SlowEndpoint)

	// Movies endpoints
	mux.HandleFunc("GET /movies", handler.GetAllMovies)
	mux.HandleFunc("POST /movies", handler.CreateMovie)
	mux.HandleFunc("GET /movies/{id}", handler.GetMovieByID)
	mux.HandleFunc("PUT /movies/{id}", handler.UpdateMovie)
	mux.HandleFunc("DELETE /movies/{id}", handler.DeleteMovie)

	// Cinemas endpoints
	mux.HandleFunc("GET /cinemas", handler.GetAllCinemas)
	mux.HandleFunc("POST /cinemas", handler.CreateCinema)
	mux.HandleFunc("GET /cinemas/{id}", handler.GetCinemaByID)

	// Halls endpoints
	mux.HandleFunc("GET /cinemas/{cinemaId}/halls", handler.GetHallsByCinema)

	// Sessions endpoints
	mux.HandleFunc("GET /sessions", handler.GetAllSessions)
	mux.HandleFunc("POST /sessions", handler.CreateSession)
	mux.HandleFunc("GET /sessions/{id}", handler.GetSessionByID)
	mux.HandleFunc("GET /movies/{movieId}/sessions", handler.GetSessionsByMovie)

	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      mux,
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
