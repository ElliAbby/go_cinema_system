package usecase

import (
	"context"
	"log"

	"github.com/ElliAbby/go_cinema_system/internal/orderService"
)

type useCase struct {
	repo orderService.Repository
}

func New(repo orderService.Repository) *useCase {
	return &useCase{repo: repo}
}

// ProcessPayment обрабатывает платёжный запрос: завершает покупку и выпускает билеты
func (uc *useCase) ProcessPayment(ctx context.Context, userID int, bookingID string) (*orderService.Booking, []orderService.Ticket, error) {
	booking, tickets, err := uc.repo.PurchaseBooking(ctx, userID, bookingID)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Payment processed: booking=%s, user=%d, tickets=%d, amount=%.2f", booking.ID, userID, len(tickets), booking.TotalPrice)
	return booking, tickets, nil
}
