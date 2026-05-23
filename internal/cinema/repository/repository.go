package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	cinemaService "github.com/ElliAbby/go_cinema_system/internal/cinema"

)

type repo struct {
 	postgresDB *sqlx.DB
}

func New(postgresDB *sqlx.DB) *repo {
	return &repo{postgresDB: postgresDB}
}

// Работа с фильмами
func (r *repo) GetAllMovies(ctx context.Context) ([]cinemaService.Movie, error) {
	var movies []cinemaService.Movie
	query := `SELECT id, title, duration, rating, description FROM movies ORDER BY id`
	err := r.postgresDB.SelectContext(ctx, &movies, query)
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return movies, nil
}

func (r *repo) GetMovieByID(ctx context.Context, id int) (*cinemaService.Movie, error) {
	var movie cinemaService.Movie
	query := `SELECT id, title, duration, rating, description FROM movies WHERE id = $1`
	err := r.postgresDB.GetContext(ctx, &movie, query, id)
	if err == sql.ErrNoRows {
		return nil, cinemaService.NewNotFoundError("movie")
	}
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return &movie, nil
}

func (r *repo) CreateMovie(ctx context.Context, movie *cinemaService.Movie) (int, error) {
	var id int
	query := `INSERT INTO movies (title, duration, rating, description) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.postgresDB.QueryRowContext(ctx, query, movie.Title, movie.Duration, movie.Rating, movie.Description).Scan(&id)
	if err != nil {
		return 0, cinemaService.NewDatabaseError(err.Error())
	}
	return id, nil
}

func (r *repo) UpdateMovie(ctx context.Context, movie *cinemaService.Movie) error {
	query := `UPDATE movies SET title = $1, duration = $2, rating = $3, description = $4 WHERE id = $5`
	result, err := r.postgresDB.ExecContext(ctx, query, movie.Title, movie.Duration, movie.Rating, movie.Description, movie.ID)
	if err != nil {
		return cinemaService.NewDatabaseError(err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return cinemaService.NewDatabaseError(err.Error())
	}
	if rowsAffected == 0 {
		return cinemaService.NewNotFoundError("movie")
	}
	return nil
}

func (r *repo) DeleteMovie(ctx context.Context, id int) error {
	query := `DELETE FROM movies WHERE id = $1`
	result, err := r.postgresDB.ExecContext(ctx, query, id)
	if err != nil {
		return cinemaService.NewDatabaseError(err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return cinemaService.NewDatabaseError(err.Error())
	}
	if rowsAffected == 0 {
		return cinemaService.NewNotFoundError("movie")
	}
	return nil
}

// Работа с кинотеатрами
func (r *repo) GetAllCinemas(ctx context.Context) ([]cinemaService.Cinema, error) {
	var cinemas []cinemaService.Cinema
	query := `SELECT id, name, address FROM cinemas ORDER BY id`
	err := r.postgresDB.SelectContext(ctx, &cinemas, query)
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return cinemas, nil
}

func (r *repo) GetCinemaByID(ctx context.Context, id int) (*cinemaService.Cinema, error) {
	var cinema cinemaService.Cinema
	query := `SELECT id, name, address FROM cinemas WHERE id = $1`
	err := r.postgresDB.GetContext(ctx, &cinema, query, id)
	if err == sql.ErrNoRows {
		return nil, cinemaService.NewNotFoundError("cinema")
	}
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return &cinema, nil
}

func (r *repo) CreateCinema(ctx context.Context, cinema *cinemaService.Cinema) (int, error) {
	var id int
	query := `INSERT INTO cinemas (name, address) VALUES ($1, $2) RETURNING id`
	err := r.postgresDB.QueryRowContext(ctx, query, cinema.Name, cinema.Address).Scan(&id)
	if err != nil {
		return 0, cinemaService.NewDatabaseError(err.Error())
	}
	return id, nil
}

// Работа с залами
func (r *repo) GetHallsByCinema(ctx context.Context, cinemaID int) ([]cinemaService.Hall, error) {
	var halls []cinemaService.Hall
	query := `SELECT id, cinema_id, name, hall_type FROM halls WHERE cinema_id = $1 ORDER BY id`
	err := r.postgresDB.SelectContext(ctx, &halls, query, cinemaID)
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return halls, nil
}

func (r *repo) GetHallByID(ctx context.Context, id int) (*cinemaService.Hall, error) {
	var hall cinemaService.Hall
	query := `SELECT id, cinema_id, name, hall_type FROM halls WHERE id = $1`
	err := r.postgresDB.GetContext(ctx, &hall, query, id)
	if err == sql.ErrNoRows {
		return nil, cinemaService.NewNotFoundError("hall")
	}
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return &hall, nil
}

func (r *repo) CreateHall(ctx context.Context, hall *cinemaService.Hall) (int, error) {
	var id int
	query := `INSERT INTO halls (cinema_id, name, hall_type) VALUES ($1, $2, $3) RETURNING id`
	err := r.postgresDB.QueryRowContext(ctx, query, hall.CinemaID, hall.Name, hall.HallType).Scan(&id)
	if err != nil {
		return 0, cinemaService.NewDatabaseError(err.Error())
	}
	return id, nil
}

// Работа с сесссиями
func (r *repo) GetAllSessions(ctx context.Context) ([]cinemaService.Session, error) {
	var sessions []cinemaService.Session
	query := `SELECT id, movie_id, hall_id, start_time, price_base FROM sessions ORDER BY start_time`
	err := r.postgresDB.SelectContext(ctx, &sessions, query)
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return sessions, nil
}

func (r *repo) GetSessionByID(ctx context.Context, id int) (*cinemaService.Session, error) {
	var session cinemaService.Session
	query := `SELECT id, movie_id, hall_id, start_time, price_base FROM sessions WHERE id = $1`
	err := r.postgresDB.GetContext(ctx, &session, query, id)
	if err == sql.ErrNoRows {
		return nil, cinemaService.NewNotFoundError("session")
	}
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return &session, nil
}

func (r *repo) GetSessionsByMovie(ctx context.Context, movieID int) ([]cinemaService.Session, error) {
	var sessions []cinemaService.Session
	query := `SELECT id, movie_id, hall_id, start_time, price_base FROM sessions WHERE movie_id = $1 ORDER BY start_time`
	err := r.postgresDB.SelectContext(ctx, &sessions, query, movieID)
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return sessions, nil
}

func (r *repo) GetSessionsByHall(ctx context.Context, hallID int) ([]cinemaService.Session, error) {
	var sessions []cinemaService.Session
	query := `SELECT id, movie_id, hall_id, start_time, price_base FROM sessions WHERE hall_id = $1 ORDER BY start_time`
	err := r.postgresDB.SelectContext(ctx, &sessions, query, hallID)
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return sessions, nil
}

func (r *repo) CreateSession(ctx context.Context, session *cinemaService.Session) (int, error) {
	var id int
	query := `INSERT INTO sessions (movie_id, hall_id, start_time, price_base) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.postgresDB.QueryRowContext(ctx, query, session.MovieID, session.HallID, session.StartTime, session.PriceBase).Scan(&id)
	if err != nil {
		return 0, cinemaService.NewDatabaseError(err.Error())
	}
	return id, nil
}

// Бронирование и покупка билетов
func (r *repo) CreateBooking(ctx context.Context, userID int, req *cinemaService.CreateBookingRequest) (*cinemaService.Booking, error) {
	tx, err := r.postgresDB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var session cinemaService.Session
	sessionQuery := `SELECT id, movie_id, hall_id, start_time, price_base FROM sessions WHERE id = $1 FOR UPDATE`
	if err := tx.GetContext(ctx, &session, sessionQuery, req.SessionID); err != nil {
		if err == sql.ErrNoRows {
			return nil, cinemaService.NewNotFoundError("session")
		}
		return nil, cinemaService.NewDatabaseError(err.Error())
	}

	seatsQuery := `SELECT id, hall_id, row_number, seat_number, seat_type FROM seats WHERE id = ANY($1)`
	var seats []cinemaService.Seat
	if err := tx.SelectContext(ctx, &seats, seatsQuery, pq.Array(req.SeatIDs)); err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	if len(seats) != len(req.SeatIDs) {
		return nil, cinemaService.NewValidationError("seat ids")
	}
	for _, seat := range seats {
		if seat.HallID != session.HallID {
			return nil, cinemaService.NewConflictError("one or more seats do not belong to this session hall")
		}
	}

	var occupiedSeatIDs []int
	occupiedQuery := `
	SELECT DISTINCT seat_id FROM (
		SELECT r.seat_id
		FROM reservations r
		JOIN bookings b ON b.id = r.booking_id
		WHERE r.session_id = $1 AND b.status IN ('pending', 'paid') AND r.seat_id = ANY($2)
		UNION
		SELECT t.seat_id
		FROM tickets t
		JOIN bookings b ON b.id = t.booking_id
		WHERE t.session_id = $1 AND b.status IN ('pending', 'paid') AND t.seat_id = ANY($2)
	) occupied
	ORDER BY seat_id
	`
	if err := tx.SelectContext(ctx, &occupiedSeatIDs, occupiedQuery, req.SessionID, pq.Array(req.SeatIDs)); err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	if len(occupiedSeatIDs) > 0 {
		return nil, cinemaService.NewConflictError("one or more seats are already reserved")
	}

	deleteExpiredQuery := `DELETE FROM reservations WHERE session_id = $1 AND seat_id = ANY($2) AND locked_until < NOW()`
	if _, err := tx.ExecContext(ctx, deleteExpiredQuery, req.SessionID, pq.Array(req.SeatIDs)); err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}

	totalPrice := session.PriceBase * float64(len(req.SeatIDs))
	var booking cinemaService.Booking
	bookingQuery := `INSERT INTO bookings (user_id, total_price, status, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id, created_at`
	if err := tx.QueryRowContext(ctx, bookingQuery, userID, totalPrice, "pending").Scan(&booking.ID, &booking.CreatedAt); err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}

	lockedUntil := time.Now().Add(15 * time.Minute)
	reservationQuery := `INSERT INTO reservations (seat_id, session_id, user_id, booking_id, locked_until) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (seat_id, session_id) DO NOTHING RETURNING seat_id`
	for _, seatID := range req.SeatIDs {
		var insertedSeatID int
		if err := tx.QueryRowContext(ctx, reservationQuery, seatID, req.SessionID, userID, booking.ID, lockedUntil).Scan(&insertedSeatID); err != nil {
			if err == sql.ErrNoRows {
				return nil, cinemaService.NewConflictError("one or more seats are already reserved")
			}
			return nil, cinemaService.NewDatabaseError(err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}

	booking.UserID = userID
	booking.TotalPrice = totalPrice
	booking.Status = "pending"
	booking.SeatIDs = append([]int(nil), req.SeatIDs...)
	booking.ExpiresAt = lockedUntil
	return &booking, nil
}

func (r *repo) GetBookingByID(ctx context.Context, userID int, bookingID string) (*cinemaService.Booking, error) {
	var booking cinemaService.Booking
	query := `SELECT id, user_id, total_price, status, created_at FROM bookings WHERE id = $1 AND user_id = $2`
	if err := r.postgresDB.GetContext(ctx, &booking, query, bookingID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, cinemaService.NewNotFoundError("booking")
		}
		return nil, cinemaService.NewDatabaseError(err.Error())
	}

	var sessionID int
    sessionQuery := `
    SELECT COALESCE(
        (SELECT r.session_id FROM reservations r WHERE r.booking_id = $1 LIMIT 1),
        (SELECT t.session_id FROM tickets t WHERE t.booking_id = $1 LIMIT 1),
        0
    ) AS session_id`
    
    if err := r.postgresDB.GetContext(ctx, &sessionID, sessionQuery, bookingID); err != nil {
        return nil, cinemaService.NewNotFoundError("session")
    }
    booking.SessionID = sessionID

	seatIDsQuery := `
	SELECT DISTINCT seat_id FROM (
		SELECT seat_id FROM reservations WHERE booking_id = $1
		UNION
		SELECT seat_id FROM tickets WHERE booking_id = $1
	) s
	ORDER BY seat_id
	`
	if err := r.postgresDB.SelectContext(ctx, &booking.SeatIDs, seatIDsQuery, bookingID); err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}

	return &booking, nil
}

func (r *repo) UpdateBookingStatus(ctx context.Context, bookingID string, status string) error {
	result, err := r.postgresDB.ExecContext(ctx, `UPDATE bookings SET status = $1 WHERE id = $2`, status, bookingID)
	if err != nil {
		return cinemaService.NewDatabaseError(err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return cinemaService.NewDatabaseError(err.Error())
	}
	if rowsAffected == 0 {
		return cinemaService.NewNotFoundError("booking")
	}
	return nil
}

func (r *repo) CancelBooking(ctx context.Context, userID int, bookingID string) error {
	tx, err := r.postgresDB.BeginTxx(ctx, nil)
	if err != nil {
		return cinemaService.NewDatabaseError(err.Error())
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var status string
	query := `SELECT status FROM bookings WHERE id = $1 AND user_id = $2 FOR UPDATE`
	if err := tx.GetContext(ctx, &status, query, bookingID, userID); err != nil {
		if err == sql.ErrNoRows {
			return cinemaService.NewNotFoundError("booking")
		}
		return cinemaService.NewDatabaseError(err.Error())
	}

	if status != "pending" {
		return cinemaService.NewConflictError("only pending bookings can be cancelled")
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM reservations WHERE booking_id = $1`, bookingID); err != nil {
		return cinemaService.NewDatabaseError(err.Error())
	}

	if _, err := tx.ExecContext(ctx, `UPDATE tickets SET status = 'refunded' WHERE booking_id = $1`, bookingID); err != nil {
		return cinemaService.NewDatabaseError(err.Error())
	}

	result, err := tx.ExecContext(ctx, `UPDATE bookings SET status = 'cancelled' WHERE id = $1 AND user_id = $2`, bookingID, userID)
	if err != nil {
		return cinemaService.NewDatabaseError(err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return cinemaService.NewDatabaseError(err.Error())
	}
	if rowsAffected == 0 {
		return cinemaService.NewNotFoundError("booking")
	}

	if err := tx.Commit(); err != nil {
		return cinemaService.NewDatabaseError(err.Error())
	}

	return nil
}

func (r *repo) GetAllMyBookings(ctx context.Context, userID int) ([]cinemaService.Booking, error) {
	var bookings []cinemaService.Booking
	query := `SELECT id, user_id, total_price, status, created_at FROM bookings WHERE user_id = $1 ORDER BY created_at DESC`
	err := r.postgresDB.SelectContext(ctx, &bookings, query, userID)
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return bookings, nil
}

// Работа с местами
func (r *repo) GetSeatsByHall(ctx context.Context, hallID int) ([]cinemaService.Seat, error) {
	var seats []cinemaService.Seat
	query := `SELECT id, hall_id, row_number, seat_number, seat_type FROM seats WHERE hall_id = $1 ORDER BY row_number, seat_number`
	err := r.postgresDB.SelectContext(ctx, &seats, query, hallID)
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return seats, nil
}

func (r *repo) GetReservedSeatIDsBySession(ctx context.Context, sessionID int) ([]int, error) {
	var seatIDs []int
	query := `
	SELECT DISTINCT seat_id FROM (
		SELECT r.seat_id
		FROM reservations r
		JOIN bookings b ON b.id = r.booking_id
		WHERE r.session_id = $1 AND b.status IN ('pending', 'paid')
		UNION
		SELECT t.seat_id
		FROM tickets t
		JOIN bookings b ON b.id = t.booking_id
		WHERE t.session_id = $1 AND b.status IN ('pending', 'paid')
	) t
	ORDER BY seat_id
	`
	if err := r.postgresDB.SelectContext(ctx, &seatIDs, query, sessionID); err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return seatIDs, nil
}

func (r *repo) GetSeatByID(ctx context.Context, id int) (*cinemaService.Seat, error) {
	var seat cinemaService.Seat
	query := `SELECT id, hall_id, row_number, seat_number, seat_type FROM seats WHERE id = $1`
	err := r.postgresDB.GetContext(ctx, &seat, query, id)
	if err == sql.ErrNoRows {
		return nil, cinemaService.NewNotFoundError("seat")
	}
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return &seat, nil
}

// Auth методы
func (r *repo) CreateUser(ctx context.Context, user *cinemaService.User) (int, error) {
	var id int
	query := `INSERT INTO users (email, password_hash, phone, is_active, created_at, updated_at) 
              VALUES ($1, $2, $3, true, NOW(), NOW()) RETURNING id`
	err := r.postgresDB.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.Phone).Scan(&id)
	if err != nil {
		return 0, cinemaService.NewDatabaseError(err.Error())
	}
	return id, nil
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*cinemaService.User, error) {
	var user cinemaService.User
	query := `SELECT id, email, password_hash, phone, is_active, created_at, updated_at FROM users WHERE email = $1`
	err := r.postgresDB.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, cinemaService.NewNotFoundError("user")
	}
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return &user, nil
}

func (r *repo) GetUserByID(ctx context.Context, id int) (*cinemaService.User, error) {
	var user cinemaService.User
	query := `SELECT id, email, password_hash, phone, is_active, created_at, updated_at FROM users WHERE id = $1`
	err := r.postgresDB.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, cinemaService.NewNotFoundError("user")
	}
	if err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return &user, nil
}

func (r *repo) GetAllUsers(ctx context.Context) ([]cinemaService.User, error) {
	var users []cinemaService.User
	query := `SELECT id, email, password_hash, phone, is_active, created_at, updated_at FROM users ORDER BY id`
	if err := r.postgresDB.SelectContext(ctx, &users, query); err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return users, nil
}

// Работа с билетами
func (r *repo) GetAllTickets(ctx context.Context, userID int) ([]cinemaService.Ticket, error) {
	var tickets []cinemaService.Ticket
	query := `SELECT t.id, t.session_id, s.start_time AS session_start_time, t.seat_id, t.booking_id, t.status FROM tickets t JOIN bookings b ON b.id = t.booking_id JOIN sessions s ON s.id = t.session_id WHERE b.user_id = $1 ORDER BY t.id`
	if err := r.postgresDB.SelectContext(ctx, &tickets, query, userID); err != nil {
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return tickets, nil
}

func (r *repo) GetTicketByID(ctx context.Context, userID int, ticketID int) (*cinemaService.Ticket, error) {
	var ticket cinemaService.Ticket
	query := `SELECT t.id, t.session_id, s.start_time AS session_start_time, t.seat_id, t.booking_id, t.status FROM tickets t JOIN bookings b ON b.id = t.booking_id JOIN sessions s ON s.id = t.session_id WHERE t.id = $1 AND b.user_id = $2`
	if err := r.postgresDB.GetContext(ctx, &ticket, query, ticketID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, cinemaService.NewNotFoundError("ticket")
		}
		return nil, cinemaService.NewDatabaseError(err.Error())
	}
	return &ticket, nil
}

// Тестовые методы
func (r *repo) GetTestMessage(ctx context.Context) (string, error) {
	return "Hello from Database!", nil
}

func (r *repo) GetSlowMessage(ctx context.Context) (string, error) {
	return "Slooooow text from Database", nil
}
