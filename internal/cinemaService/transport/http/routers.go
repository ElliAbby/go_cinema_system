package http

import (
    "net/http"

    "github.com/go-chi/chi/v5"

)

func RegisterRouters(h *handler) http.Handler {
    r := chi.NewRouter()

	r.Use(metricsMiddleware)

    // Root
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"message":"Добро пожаловать в кинотеатр API!"}`))
    })

    // Health
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"healthy"}`))
    })

	// Metrics
	r.Handle("/metrics", metricsHandler())

    // Auth
    r.Post("/auth/register", h.Register)
    r.Post("/auth/login", h.Login)
    // TODO: добавить рефреш токен
    // r.Post("/auth/refresh", h.Refresh) 

    // Test
    r.Get("/test", h.TestEndpoint)
    r.Get("/slow", h.SlowEndpoint)

    // Movies
    r.Route("/movies", func(r chi.Router) {
        r.Get("/", h.GetAllMovies)
        r.Post("/", func(w http.ResponseWriter, r *http.Request) {
            JWTMiddleware(http.HandlerFunc(h.CreateMovie)).ServeHTTP(w, r)
        })
        r.Get("/{id}", h.GetMovieByID)
        r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
            JWTMiddleware(http.HandlerFunc(h.UpdateMovie)).ServeHTTP(w, r)
        })
        r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
            JWTMiddleware(http.HandlerFunc(h.DeleteMovie)).ServeHTTP(w, r)
        })
        r.Get("/{id}/sessions", h.GetSessionsByMovie)
    })

    // Cinemas
    r.Route("/cinemas", func(r chi.Router) {
        r.Get("/", h.GetAllCinemas)
        r.Post("/", func(w http.ResponseWriter, r *http.Request) {
            JWTMiddleware(http.HandlerFunc(h.CreateCinema)).ServeHTTP(w, r)
        })
        r.Get("/{id}", h.GetCinemaByID)
        r.Get("/{id}/halls", h.GetHallsByCinema)
        // TODO: добавить эдндпоинты для залов
    })

    // Sessions
    r.Route("/sessions", func(r chi.Router) {
        r.Get("/", h.GetAllSessions)
        r.Post("/", func(w http.ResponseWriter, r *http.Request) {
            JWTMiddleware(http.HandlerFunc(h.CreateSession)).ServeHTTP(w, r)
        })
        r.Get("/{id}", h.GetSessionByID)
        // TODO: добавить список доступных мест
    })

    // Bookings
    r.Route("/bookings", func(r chi.Router) {
        r.Get("/", func(w http.ResponseWriter, r *http.Request) {
            JWTMiddleware(http.HandlerFunc(h.GetAllMyBookings)).ServeHTTP(w, r)
        })
        r.Post("/", func(w http.ResponseWriter, r *http.Request) {
            JWTMiddleware(http.HandlerFunc(h.CreateBooking)).ServeHTTP(w, r)
        })
        r.Post("/{id}/purchase", func(w http.ResponseWriter, r *http.Request) {
            JWTMiddleware(http.HandlerFunc(h.PurchaseBooking)).ServeHTTP(w, r)
        })
        // TODO: добавить удаление брони - до опалты
        // r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
        //     JWTMiddleware(http.HandlerFunc(h.DeleteBooking)).ServeHTTP(w, r)
        // })
    })

    // Users
    r.Route("/users", func(r chi.Router) {
        r.Get("/", h.ListUsers)
        r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
            JWTMiddleware(http.HandlerFunc(h.GetMe)).ServeHTTP(w, r)
        })
        // TODO: добавить обновление пользователя
        // r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
        //     JWTMiddleware(http.HandlerFunc(h.UpdateMe)).ServeHTTP(w, r)
        // })
    })

    // Tickets
    r.Route("/tickets", func(r chi.Router) {
        r.Get("/", func(w http.ResponseWriter, r *http.Request) {
            JWTMiddleware(http.HandlerFunc(h.ListTickets)).ServeHTTP(w, r)
        })
        r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
            JWTMiddleware(http.HandlerFunc(h.GetTicketByID)).ServeHTTP(w, r)
        })
    })

    return r
}
