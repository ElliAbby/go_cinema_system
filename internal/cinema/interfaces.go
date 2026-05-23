package cinema

import "context"

// Repository определяет методы доступа к данным
type Repository interface {
	// Фильмы
	GetAllMovies(ctx context.Context) ([]Movie, error)
	GetMovieByID(ctx context.Context, id int) (*Movie, error)
	CreateMovie(ctx context.Context, movie *Movie) (int, error)
	UpdateMovie(ctx context.Context, movie *Movie) error
	DeleteMovie(ctx context.Context, id int) error

	// Кинотеатры
	GetAllCinemas(ctx context.Context) ([]Cinema, error)
	GetCinemaByID(ctx context.Context, id int) (*Cinema, error)
	CreateCinema(ctx context.Context, cinema *Cinema) (int, error)

	// Залы
	GetHallsByCinema(ctx context.Context, cinemaID int) ([]Hall, error)
	GetHallByID(ctx context.Context, id int) (*Hall, error)
	CreateHall(ctx context.Context, hall *Hall) (int, error)

	// Сессии
	GetAllSessions(ctx context.Context) ([]Session, error)
	GetSessionByID(ctx context.Context, id int) (*Session, error)
	GetSessionsByMovie(ctx context.Context, movieID int) ([]Session, error)
	GetSessionsByHall(ctx context.Context, hallID int) ([]Session, error)
	CreateSession(ctx context.Context, session *Session) (int, error)

	// Бронирования и покупка
	CreateBooking(ctx context.Context, userID int, req *CreateBookingRequest) (*Booking, error)
	GetBookingByID(ctx context.Context, userID int, bookingID string) (*Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID string, status string) error
	CancelBooking(ctx context.Context, userID int, bookingID string) (error)
	GetAllMyBookings(ctx context.Context, userID int) ([]Booking, error)

	// Reservations
	GetReservedSeatIDsBySession(ctx context.Context, sessionID int) ([]int, error)

	// Места
	GetSeatsByHall(ctx context.Context, hallID int) ([]Seat, error)
	GetSeatByID(ctx context.Context, id int) (*Seat, error)

	// Auth методы
	CreateUser(ctx context.Context, user *User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)

	// Пользователи
	GetAllUsers(ctx context.Context) ([]User, error)

	// Билеты
	GetAllTickets(ctx context.Context, userID int) ([]Ticket, error)
	GetTicketByID(ctx context.Context, userID int, ticketID int) (*Ticket, error)

	// Тестовые методы
	GetTestMessage(ctx context.Context) (string, error)
	GetSlowMessage(ctx context.Context) (string, error)
	CreateTestMessage(ctx context.Context, message string) (int, error)
}

// UseCase определяет бизнес-логику
type UseCase interface {
	// Фильмы
	GetAllMovies(ctx context.Context) ([]Movie, error)
	GetMovieByID(ctx context.Context, id int) (*Movie, error)
	CreateMovie(ctx context.Context, req *CreateMovieRequest) (int, error)
	UpdateMovie(ctx context.Context, id int, req *CreateMovieRequest) error
	DeleteMovie(ctx context.Context, id int) error

	// Кинотеатры
	GetAllCinemas(ctx context.Context) ([]Cinema, error)
	GetCinemaByID(ctx context.Context, id int) (*Cinema, error)
	CreateCinema(ctx context.Context, req *CreateCinemaRequest) (int, error)

	// Залы
	GetHallsByCinema(ctx context.Context, cinemaID int) ([]Hall, error)
	GetHallByID(ctx context.Context, id int) (*Hall, error)

	// Сессии
	GetAllSessions(ctx context.Context) ([]Session, error)
	GetSessionByID(ctx context.Context, id int) (*Session, error)
	GetSessionsByMovie(ctx context.Context, movieID int) ([]Session, error)
	CreateSession(ctx context.Context, req *CreateSessionRequest) (int, error)

	// Бронирования и покупка
	CreateBooking(ctx context.Context, userID int, req *CreateBookingRequest) (*Booking, error)
	GetBookingByID(ctx context.Context, userID int, bookingID string) (*Booking, error)
	RequestBookingPayment(ctx context.Context, userID int, bookingID string) (*Booking, error)
	CancelBooking(ctx context.Context, userID int, bookingID string) (*Booking, error)
	GetAllMyBookings(ctx context.Context, userID int) ([]Booking, error)

	// Reservations
	GetReservedSeatIDsBySession(ctx context.Context, sessionID int) ([]int, error)

	// Seats
	GetSeatsByHall(ctx context.Context, hallID int) ([]Seat, error)

	// Auth методы
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)

	// Пользователи
	ListUsers(ctx context.Context) ([]User, error)
	GetMe(ctx context.Context, userID int) (*User, error)

	// Билеты
	ListTickets(ctx context.Context, userID int) ([]Ticket, error)
	GetTicketByID(ctx context.Context, userID int, ticketID int) (*Ticket, error)

	// Тестовые методы
	GetTestMessage(ctx context.Context) (string, error)
	GetSlowMessage(ctx context.Context) (string, error)
	CreateTestMessage(ctx context.Context, message string) (int, error)
}
