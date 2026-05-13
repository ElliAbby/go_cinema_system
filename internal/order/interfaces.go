package order

import "context"

type Repository interface {
	PurchaseBooking(ctx context.Context, userID int, bookingID string) (*Booking, []Ticket, error)
}

type UseCase interface {
	ProcessPayment(ctx context.Context, userID int, bookingID string) (*Booking, []Ticket, error)
}
