package http

import (
    "net/http"

    "github.com/go-chi/chi/v5"
)

func RegisterRouters(h *handler) http.Handler {
    r := chi.NewRouter()

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

    // Auth
    r.Post("/auth/register", h.Register)
    r.Post("/auth/login", h.Login)

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
    })

    // Sessions
    r.Route("/sessions", func(r chi.Router) {
        r.Get("/", h.GetAllSessions)
        r.Post("/", func(w http.ResponseWriter, r *http.Request) {
            JWTMiddleware(http.HandlerFunc(h.CreateSession)).ServeHTTP(w, r)
        })
        r.Get("/{id}", h.GetSessionByID)
    })

    return r
}
