package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ElliAbby/go_cinema_system/internal/order"
)

type repo struct {
	postgresDB *sqlx.DB
}

func New(postgresDB *sqlx.DB) *repo {
	return &repo{postgresDB: postgresDB}
}

func (r *repo) PurchaseBooking(ctx context.Context, userID int, bookingID string) (*order.Booking, []order.Ticket, error) {
	tx, err := r.postgresDB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var booking order.Booking
	bookingQuery := `SELECT id, user_id, total_price, status FROM bookings WHERE id = $1 AND user_id = $2 FOR UPDATE`
	if err := tx.GetContext(ctx, &booking, bookingQuery, bookingID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, sql.ErrNoRows
		}
		return nil, nil, err
	}

	if booking.Status != "pending" && booking.Status != "processing" {
		return nil, nil, sql.ErrNoRows
	}

	var reservations []order.Reservation
	reservationQuery := `SELECT seat_id, session_id, user_id, booking_id, locked_until FROM reservations WHERE booking_id = $1 FOR UPDATE`
	if err := tx.SelectContext(ctx, &reservations, reservationQuery, bookingID); err != nil {
		return nil, nil, err
	}
	if len(reservations) == 0 {
		return nil, nil, sql.ErrNoRows
	}

	// временная проверка (данные о locked_until)
	var now time.Time
	if err := tx.GetContext(ctx, &now, `SELECT NOW()`); err != nil {
		return nil, nil, err
	}

	for _, res := range reservations {
		if res.LockedUntil != nil && res.LockedUntil.Before(now) {
			_, errCleanup := tx.ExecContext(ctx, `DELETE FROM reservations WHERE booking_id = $1`, bookingID)
			if errCleanup != nil {
				return nil, nil, fmt.Errorf("cleanup failed: %w", errCleanup)
			}
			_, errStatus := tx.ExecContext(ctx, `UPDATE bookings SET status = 'cancelled' WHERE id = $1`, bookingID)
			if errStatus != nil {
				return nil, nil, fmt.Errorf("status update failed: %w", errStatus)
			}
			if err := tx.Commit(); err != nil {
				return nil, nil, err
			}
			return nil, nil, order.ErrReservationExpired
		}
	}

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

	if _, err := tx.ExecContext(ctx, `UPDATE bookings SET status = $1 WHERE id = $2`, "paid", bookingID); err != nil {
		return nil, nil, err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM reservations WHERE booking_id = $1`, bookingID); err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	booking.Status = "paid"
	return &booking, tickets, nil
}
