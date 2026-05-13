package orderService

import "context"

// Repository интерфейс для работы с платежами и заказами
type Repository interface {
	// PurchaseBooking завершает покупку: создаёт билеты и обновляет статус заказа
	PurchaseBooking(ctx context.Context, userID int, bookingID string) (*Booking, []Ticket, error)
}

// UseCase интерфейс для бизнес-логики платежей
type UseCase interface {
	ProcessPayment(ctx context.Context, userID int, bookingID string) (*Booking, []Ticket, error)
}
