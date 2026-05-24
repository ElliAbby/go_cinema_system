package cinema

import "time"

// Movie представляет фильм в кинотеатре
type Movie struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Duration    int       `json:"duration" db:"duration"` // минуты
	Rating      string    `json:"rating" db:"rating"`
	Description string    `json:"description" db:"description"`
}

// Cinema представляет кинотеатр
type Cinema struct {
	ID      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Address string `json:"address" db:"address"`
}

// Hall представляет зал в кинотеатре
type Hall struct {
	ID       int    `json:"id" db:"id"`
	CinemaID int    `json:"cinema_id" db:"cinema_id"`
	Name     string `json:"name" db:"name"`
	HallType string `json:"hall_type" db:"hall_type"` // стандарт, VIP, IMAX
}

// Seat представляет место в зале
type Seat struct {
	ID        int    `json:"id" db:"id"`
	HallID    int    `json:"hall_id" db:"hall_id"`
	RowNumber int    `json:"row_number" db:"row_number"`
	SeatNum   int    `json:"seat_number" db:"seat_number"`
	SeatType  string `json:"seat_type" db:"seat_type"` // обычное, диван, для инвалидов
}

// Session представляет сеанс показа фильма
type Session struct {
	ID        int       `json:"id" db:"id"`
	MovieID   int       `json:"movie_id" db:"movie_id"`
	HallID    int       `json:"hall_id" db:"hall_id"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	PriceBase float64   `json:"price_base" db:"price_base"`
}

// User представляет пользователя
type User struct {
	ID            int       `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	Phone         string    `json:"phone" db:"phone"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Booking представляет бронирование билетов
type Booking struct {
	ID         string    `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	SessionID  int       `json:"session_id" db:"-"`
	TotalPrice float64   `json:"total_price" db:"total_price"`
	Status     string    `json:"status" db:"status"` // pending, paid, cancelled, refunded
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ExpiresAt  time.Time `json:"expires_at,omitempty" db:"-"`
	SeatIDs    []int     `json:"seat_ids,omitempty" db:"-"`
}

// Ticket представляет билет
type Ticket struct {
	ID               string    `json:"id" db:"id"`
	SessionID        int       `json:"session_id" db:"session_id"`
	SessionStartTime time.Time `json:"session_start_time,omitempty" db:"session_start_time"`
	RowNumber        int       `json:"row_number,omitempty" db:"row_number"`
	SeatNumber       int       `json:"seat_number,omitempty" db:"seat_number"`
	SeatID           int       `json:"seat_id" db:"seat_id"`
	BookingID        string    `json:"booking_id" db:"booking_id"`
	Status           string    `json:"status" db:"status"` // active, used, deactivated
}

// Reservation представляет временную фиксацию места перед покупкой
type Reservation struct {
	SeatID      int       `json:"seat_id" db:"seat_id"`
	SessionID   int       `json:"session_id" db:"session_id"`
	UserID      int       `json:"user_id" db:"user_id"`
	BookingID   string    `json:"booking_id" db:"booking_id"`
	LockedUntil time.Time `json:"locked_until" db:"locked_until"`
}

// CreateMovieRequest для создания фильма
type CreateMovieRequest struct {
	Title       string `json:"title" binding:"required"`
	Duration    int    `json:"duration" binding:"required,min=1"`
	Rating      string `json:"rating"`
	Description string `json:"description"`
}

// CreateCinemaRequest для создания кинотеатра
type CreateCinemaRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

// CreateSessionRequest для создания сеанса
type CreateSessionRequest struct {
	MovieID   int       `json:"movie_id" binding:"required"`
	HallID    int       `json:"hall_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	PriceBase float64   `json:"price_base" binding:"required,min=0"`
}

// CreateBookingRequest для бронирования мест на сеанс
type CreateBookingRequest struct {
	SessionID int   `json:"session_id" binding:"required"`
	SeatIDs   []int `json:"seat_ids" binding:"required,min=1"`
}

// Модели для работы с авторизацией и аутентификацией
type RegisterRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Phone string `json:"phone"`
}

// LoginRequest для аутентификации
type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse для ответа на запрос аутентификации
type AuthResponse struct {
	Token string `json:"token"`
	UserID  int   `json:"user_id"`
	Email string `json:"email"`
	ExpiresAt int64 `json:"expires_at"`
}