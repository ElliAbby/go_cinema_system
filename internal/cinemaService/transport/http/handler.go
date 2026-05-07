package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/ElliAbby/go_cinema_system/internal/cinemaService"

)

type handler struct {
	uc cinemaService.UseCase
}

func New(uc cinemaService.UseCase) *handler {
	return &handler{uc: uc}
}

// Дополнительные функции для ответа и извлечения параметров
func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, err interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err,
	})
}

func getIDFromURL(r *http.Request, paramName string) (int, error) {
	idStr := chi.URLParam(r, paramName)
	return strconv.Atoi(idStr)
}

// Эндпоинты для фильмов
func (h *handler) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.uc.GetAllMovies(r.Context())
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, movies, http.StatusOK)
}

func (h *handler) GetMovieByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r, "id")
	if err != nil {
		respondError(w, cinemaService.NewValidationError("id"), 400)
		return
	}

	movie, err := h.uc.GetMovieByID(r.Context(), id)
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, movie, http.StatusOK)
}

func (h *handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var req cinemaService.CreateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, cinemaService.ErrInvalidInput, 400)
		return
	}

	id, err := h.uc.CreateMovie(r.Context(), &req)
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, map[string]int{"id": id}, http.StatusCreated)
}

func (h *handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r, "id")
	if err != nil {
		respondError(w, cinemaService.NewValidationError("id"), 400)
		return
	}

	var req cinemaService.CreateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, cinemaService.ErrInvalidInput, 400)
		return
	}

	err = h.uc.UpdateMovie(r.Context(), id, &req)
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, map[string]string{"message": "movie updated"}, http.StatusOK)
}

func (h *handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r, "id")
	if err != nil {
		respondError(w, cinemaService.NewValidationError("id"), 400)
		return
	}

	err = h.uc.DeleteMovie(r.Context(), id)
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, map[string]string{"message": "movie deleted"}, http.StatusOK)
}

// Эндпоинты для кинотеатров
func (h *handler) GetAllCinemas(w http.ResponseWriter, r *http.Request) {
	cinemas, err := h.uc.GetAllCinemas(r.Context())
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, cinemas, http.StatusOK)
}

func (h *handler) GetCinemaByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r, "id")
	if err != nil {
		respondError(w, cinemaService.NewValidationError("id"), 400)
		return
	}

	cinema, err := h.uc.GetCinemaByID(r.Context(), id)
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, cinema, http.StatusOK)
}

func (h *handler) CreateCinema(w http.ResponseWriter, r *http.Request) {
	var req cinemaService.CreateCinemaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, cinemaService.ErrInvalidInput, 400)
		return
	}

	id, err := h.uc.CreateCinema(r.Context(), &req)
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, map[string]int{"id": id}, http.StatusCreated)
}

// Эндпоинты для залов
func (h *handler) GetHallsByCinema(w http.ResponseWriter, r *http.Request) {
	cinemaID, err := getIDFromURL(r, "cinemaId")
	if err != nil {
		respondError(w, cinemaService.NewValidationError("cinemaId"), 400)
		return
	}

	halls, err := h.uc.GetHallsByCinema(r.Context(), cinemaID)
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, halls, http.StatusOK)
}

// Эндпоинты для сессий
func (h *handler) GetAllSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.uc.GetAllSessions(r.Context())
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, sessions, http.StatusOK)
}

func (h *handler) GetSessionByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r, "id")
	if err != nil {
		respondError(w, cinemaService.NewValidationError("id"), 400)
		return
	}

	session, err := h.uc.GetSessionByID(r.Context(), id)
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, session, http.StatusOK)
}

func (h *handler) GetSessionsByMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := getIDFromURL(r, "movieId")
	if err != nil {
		respondError(w, cinemaService.NewValidationError("movieId"), 400)
		return
	}

	sessions, err := h.uc.GetSessionsByMovie(r.Context(), movieID)
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, sessions, http.StatusOK)
}

func (h *handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req cinemaService.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, cinemaService.ErrInvalidInput, 400)
		return
	}

	id, err := h.uc.CreateSession(r.Context(), &req)
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, map[string]int{"id": id}, http.StatusCreated)
}

// Auth middleware
func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var req cinemaService.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, cinemaService.ErrInvalidInput, 400)
		return
	}

	resp, err := h.uc.Register(r.Context(), &req)
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}
	respondJSON(w, resp, http.StatusOK)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req cinemaService.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, cinemaService.ErrInvalidInput, 400)
		return
	}

	resp, err := h.uc.Login(r.Context(), &req)
	if err != nil {
		if appErr, ok := err.(cinemaService.AppError); ok {
			respondError(w, appErr, appErr.Status)
		} else {
			respondError(w, cinemaService.ErrInternal, 500)
		}
		return
	}

	respondJSON(w, resp, http.StatusOK)
}

// Тестовые эндпоинты
func (h *handler) TestEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос дошел до хэндлера /test!")
	msg, err := h.uc.GetTestMessage(r.Context())
	if err != nil {
		respondError(w, cinemaService.ErrInternal, 500)
		return
	}
	respondJSON(w, map[string]string{"message": msg}, http.StatusOK)
}

func (h *handler) SlowEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос дошел до медленного хэндлера /slow!")

	select {
	case <-time.After(5 * time.Second):
	case <-r.Context().Done():
		log.Println("Запрос прерван до завершения паузы")
		return
	}

	msg, err := h.uc.GetSlowMessage(r.Context())
	if err != nil {
		respondError(w, cinemaService.ErrInternal, 500)
		return
	}
	respondJSON(w, map[string]string{"message": msg}, http.StatusOK)
}
