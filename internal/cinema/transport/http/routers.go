package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	authjwt "github.com/ElliAbby/go_cinema_system/internal/platform/auth/jwt"
)

func RegisterRouters(h *handler, tokenManager *authjwt.Manager) http.Handler {
	r := chi.NewRouter()

	r.Use(corsMiddleware)
	r.Use(metricsMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Добро пожаловать в кинотеатр API!"}`))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy"}`))
	})

	r.Handle("/metrics", metricsHandler())

	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)

	r.Get("/test", h.TestEndpoint)
	r.Get("/slow", h.SlowEndpoint)
	r.Post("/dbtest", h.DBTestEndpoint)

	r.Route("/movies", func(r chi.Router) {
		r.Get("/", h.GetAllMovies)
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.CreateMovie)).ServeHTTP(w, r)
		})
		r.Get("/{id}", h.GetMovieByID)
		r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.UpdateMovie)).ServeHTTP(w, r)
		})
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.DeleteMovie)).ServeHTTP(w, r)
		})
		r.Get("/{id}/sessions", h.GetSessionsByMovie)
	})

	r.Route("/cinemas", func(r chi.Router) {
		r.Get("/", h.GetAllCinemas)
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.CreateCinema)).ServeHTTP(w, r)
		})
		r.Get("/{id}", h.GetCinemaByID)
		r.Get("/{id}/halls", h.GetHallsByCinema)
	})

	// Seats and halls
	r.Get("/halls/{id}", h.GetHallByID)
	r.Get("/halls/{id}/seats", h.GetSeatsByHall)

	r.Route("/sessions", func(r chi.Router) {
		r.Get("/", h.GetAllSessions)
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.CreateSession)).ServeHTTP(w, r)
		})
		r.Get("/{id}", h.GetSessionByID)
		r.Get("/{id}/reserved", h.GetReservedSeats)
	})

	r.Route("/bookings", func(r chi.Router) {
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.GetBookingByID)).ServeHTTP(w, r)
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.GetAllMyBookings)).ServeHTTP(w, r)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.CreateBooking)).ServeHTTP(w, r)
		})
		r.Post("/{id}/purchase", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.PurchaseBooking)).ServeHTTP(w, r)
		})
		r.Post("/{id}/cancel", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.CancelBooking)).ServeHTTP(w, r)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.ListUsers)
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.GetMe)).ServeHTTP(w, r)
		})
	})

	r.Route("/tickets", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.ListTickets)).ServeHTTP(w, r)
		})
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			JWTMiddleware(tokenManager, http.HandlerFunc(h.GetTicketByID)).ServeHTTP(w, r)
		})
	})

	return r
}
