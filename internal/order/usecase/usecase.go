package usecase

import (
	"context"
	"log"

	"github.com/ElliAbby/go_cinema_system/internal/order"

)

type useCase struct {
	repo order.Repository
}

func New(repo order.Repository) *useCase {
	return &useCase{repo: repo}
}

func (uc *useCase) ProcessPayment(ctx context.Context, userID int, bookingID string) (*order.Booking, []order.Ticket, error) {
	booking, tickets, err := uc.repo.PurchaseBooking(ctx, userID, bookingID)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Payment processed: booking=%s, user=%d, tickets=%d, amount=%.2f", booking.ID, userID, len(tickets), booking.TotalPrice)
	return booking, tickets, nil
}
