package http

import (
	"log"
	"time"
	"net/http"

	"github.com/ElliAbby/go_cinema_system/internal/cinemaService"
)

type handler struct {
	uc cinemaService.UseCase
}

func New(uc cinemaService.UseCase) *handler {
	return &handler{uc: uc}
}

func (h *handler) TestEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос дошел до хэндлера /test!")
	msg, err := h.uc.GetTestMessage(r.Context())
	if err != nil{
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}

func (h *handler) SlowEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос дошел до медленного хэндлера /slow!")

	// time.Sleep(5 * time.Second) - плохо, так как полностью заблокирует текущую горутину
	select {
		case <-time.After(5 * time.Second):  // успешная пауза
		case <-r.Context().Done():  // немедленно прекращаем работу
			log.Println("Запрос прерван до завершения паузы")
			return
	}

	msg, err := h.uc.GetSlowMessage(r.Context())
	if err != nil{
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}
