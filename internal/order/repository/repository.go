package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/ElliAbby/go_cinema_system/internal/order"
)

type repo struct {
	postgresDB *sqlx.DB
}

func New(postgresDB *sqlx.DB) *repo {
	return &repo{postgresDB: postgresDB}
}

// PurchaseBooking завершает покупку: создаёт билеты, обновляет статус на "paid" и удаляет резервации
func (r *repo) PurchaseBooking(ctx context.Context, userID int, bookingID string) (*order.Booking, []order.Ticket, error) {
	tx, err := r.postgresDB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Получаем заказ с блокировкой для предотвращения race condition
	var booking order.Booking
	bookingQuery := `SELECT id, user_id, total_price, status FROM bookings WHERE id = $1 AND user_id = $2 FOR UPDATE`
	if err := tx.GetContext(ctx, &booking, bookingQuery, bookingID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, sql.ErrNoRows
		}
		return nil, nil, err
	}

	// Проверяем статус заказа
	if booking.Status != "pending" && booking.Status != "processing" {
		return nil, nil, sql.ErrNoRows // или собственная ошибка
	}

	// Получаем резервации мест
	var reservations []order.Reservation
	reservationQuery := `SELECT seat_id, session_id, user_id, booking_id FROM reservations WHERE booking_id = $1 FOR UPDATE`
	if err := tx.SelectContext(ctx, &reservations, reservationQuery, bookingID); err != nil {
		return nil, nil, err
	}
	if len(reservations) == 0 {
		return nil, nil, sql.ErrNoRows
	}

	// Проверяем, что резервации не истекли (данные о locked_until пока не используются)
	// В будущем можно добавить временную проверку

	// Создаём билеты
	tickets := make([]order.Ticket, 0, len(reservations))
	ticketQuery := `INSERT INTO tickets (session_id, seat_id, booking_id, status) VALUES ($1, $2, $3, $4) ON CONFLICT (session_id, seat_id) DO NOTHING RETURNING id`
	for _, reservation := range reservations {
		var ticket order.Ticket
		ticket.SessionID = reservation.SessionID
		ticket.SeatID = reservation.SeatID
		ticket.BookingID = bookingID
		ticket.Status = "active"
		if err := tx.QueryRowContext(ctx, ticketQuery, reservation.SessionID, reservation.SeatID, bookingID, ticket.Status).Scan(&ticket.ID); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil, sql.ErrNoRows
			}
			return nil, nil, err
		}
		tickets = append(tickets, ticket)
	}

	// Обновляем статус заказа на "paid"
	if _, err := tx.ExecContext(ctx, `UPDATE bookings SET status = $1 WHERE id = $2`, "paid", bookingID); err != nil {
		return nil, nil, err
	}

	// Удаляем резервации
	if _, err := tx.ExecContext(ctx, `DELETE FROM reservations WHERE booking_id = $1`, bookingID); err != nil {
		return nil, nil, err
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	booking.Status = "paid"
	return &booking, tickets, nil
}
