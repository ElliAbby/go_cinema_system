package order

import "time"

type Booking struct {
	ID         string  `json:"id" db:"id"`
	UserID     int     `json:"user_id" db:"user_id"`
	TotalPrice float64 `json:"total_price" db:"total_price"`
	Status     string  `json:"status" db:"status"`
}

type Ticket struct {
	ID        string `json:"id" db:"id"`
	SessionID int    `json:"session_id" db:"session_id"`
	SeatID    int    `json:"seat_id" db:"seat_id"`
	BookingID string `json:"booking_id" db:"booking_id"`
	Status    string `json:"status" db:"status"` // active, used, deactivated
}

type Reservation struct {
	SeatID      int        `db:"seat_id"`
	SessionID   int        `db:"session_id"`
	UserID      int        `db:"user_id"`
	BookingID   string     `db:"booking_id"`
	LockedUntil *time.Time `db:"locked_until"`
}