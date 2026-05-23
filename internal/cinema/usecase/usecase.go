package usecase

import (
	"context"
	"fmt"
	"time"

	cinema "github.com/ElliAbby/go_cinema_system/internal/cinema"
	cinemaKafka "github.com/ElliAbby/go_cinema_system/internal/platform/kafka"
	authjwt "github.com/ElliAbby/go_cinema_system/internal/platform/auth/jwt"
	"github.com/ElliAbby/go_cinema_system/internal/platform/metrics"
	"github.com/ElliAbby/go_cinema_system/internal/platform/security"

)

type useCase struct {
	repo         cinema.Repository
	publisher    PaymentPublisher
	tokenManager *authjwt.Manager
}

type PaymentPublisher interface {
	PublishPaymentRequested(ctx context.Context, event cinemaKafka.PaymentRequestedEvent) error
}

const reservationHoldDuration = 15 * time.Minute

func New(r cinema.Repository, publisher PaymentPublisher, tokenManager *authjwt.Manager) *useCase {
	return &useCase{repo: r, publisher: publisher, tokenManager: tokenManager}
}

// Фильмы
func (uc *useCase) GetAllMovies(ctx context.Context) ([]cinema.Movie, error) {
	return uc.repo.GetAllMovies(ctx)
}

func (uc *useCase) GetMovieByID(ctx context.Context, id int) (*cinema.Movie, error) {
	if id <= 0 {
		return nil, cinema.NewValidationError("movie id")
	}
	return uc.repo.GetMovieByID(ctx, id)
}

func (uc *useCase) CreateMovie(ctx context.Context, req *cinema.CreateMovieRequest) (int, error) {
	if req == nil {
		return 0, cinema.ErrInvalidInput
	}
	if req.Title == "" {
		return 0, cinema.NewValidationError("title")
	}
	if req.Duration <= 0 {
		return 0, cinema.NewValidationError("duration")
	}

	movie := &cinema.Movie{
		Title:       req.Title,
		Duration:    req.Duration,
		Rating:      req.Rating,
		Description: req.Description,
	}
	return uc.repo.CreateMovie(ctx, movie)
}

func (uc *useCase) UpdateMovie(ctx context.Context, id int, req *cinema.CreateMovieRequest) error {
	if id <= 0 {
		return cinema.NewValidationError("movie id")
	}
	if req == nil {
		return cinema.ErrInvalidInput
	}
	if req.Title == "" {
		return cinema.NewValidationError("title")
	}
	if req.Duration <= 0 {
		return cinema.NewValidationError("duration")
	}

	movie := &cinema.Movie{
		ID:          id,
		Title:       req.Title,
		Duration:    req.Duration,
		Rating:      req.Rating,
		Description: req.Description,
	}
	return uc.repo.UpdateMovie(ctx, movie)
}

func (uc *useCase) DeleteMovie(ctx context.Context, id int) error {
	if id <= 0 {
		return cinema.NewValidationError("movie id")
	}
	return uc.repo.DeleteMovie(ctx, id)
}

// Кинотеатры
func (uc *useCase) GetAllCinemas(ctx context.Context) ([]cinema.Cinema, error) {
	return uc.repo.GetAllCinemas(ctx)
}

func (uc *useCase) GetCinemaByID(ctx context.Context, id int) (*cinema.Cinema, error) {
	if id <= 0 {
		return nil, cinema.NewValidationError("cinema id")
	}
	return uc.repo.GetCinemaByID(ctx, id)
}

func (uc *useCase) CreateCinema(ctx context.Context, req *cinema.CreateCinemaRequest) (int, error) {
	if req == nil {
		return 0, cinema.ErrInvalidInput
	}
	if req.Name == "" {
		return 0, cinema.NewValidationError("name")
	}
	if req.Address == "" {
		return 0, cinema.NewValidationError("address")
	}

	cinemaModel := &cinema.Cinema{
		Name:    req.Name,
		Address: req.Address,
	}
	return uc.repo.CreateCinema(ctx, cinemaModel)
}

// Залы
func (uc *useCase) GetHallsByCinema(ctx context.Context, cinemaID int) ([]cinema.Hall, error) {
	if cinemaID <= 0 {
		return nil, cinema.NewValidationError("cinema id")
	}
	return uc.repo.GetHallsByCinema(ctx, cinemaID)
}

func (uc *useCase) GetSeatsByHall(ctx context.Context, hallID int) ([]cinema.Seat, error) {
	if hallID <= 0 {
		return nil, cinema.NewValidationError("hall id")
	}
	return uc.repo.GetSeatsByHall(ctx, hallID)
}

// Сессии
func (uc *useCase) GetAllSessions(ctx context.Context) ([]cinema.Session, error) {
	return uc.repo.GetAllSessions(ctx)
}

func (uc *useCase) GetSessionByID(ctx context.Context, id int) (*cinema.Session, error) {
	if id <= 0 {
		return nil, cinema.NewValidationError("session id")
	}
	return uc.repo.GetSessionByID(ctx, id)
}

func (uc *useCase) GetSessionsByMovie(ctx context.Context, movieID int) ([]cinema.Session, error) {
	if movieID <= 0 {
		return nil, cinema.NewValidationError("movie id")
	}
	return uc.repo.GetSessionsByMovie(ctx, movieID)
}

func (uc *useCase) CreateSession(ctx context.Context, req *cinema.CreateSessionRequest) (int, error) {
	if req == nil {
		return 0, cinema.ErrInvalidInput
	}
	if req.MovieID <= 0 {
		return 0, cinema.NewValidationError("movie id")
	}
	if req.HallID <= 0 {
		return 0, cinema.NewValidationError("hall id")
	}
	if req.PriceBase < 0 {
		return 0, cinema.NewValidationError("price base")
	}

	session := &cinema.Session{
		MovieID:   req.MovieID,
		HallID:    req.HallID,
		StartTime: req.StartTime,
		PriceBase: req.PriceBase,
	}
	return uc.repo.CreateSession(ctx, session)
}

// Бронирование и покупка билетов
func validateSeatIDs(seatIDs []int) error {
	if len(seatIDs) == 0 {
		return cinema.NewValidationError("seat ids")
	}

	seen := make(map[int]struct{}, len(seatIDs))
	for _, seatID := range seatIDs {
		if seatID <= 0 {
			return cinema.NewValidationError("seat id")
		}
		if _, exists := seen[seatID]; exists {
			return cinema.NewValidationError("seat ids")
		}
		seen[seatID] = struct{}{}
	}

	return nil
}

func (uc *useCase) CreateBooking(ctx context.Context, userID int, req *cinema.CreateBookingRequest) (*cinema.Booking, error) {
	if userID <= 0 {
		return nil, cinema.NewValidationError("user id")
	}
	if req == nil {
		return nil, cinema.ErrInvalidInput
	}
	if req.SessionID <= 0 {
		return nil, cinema.NewValidationError("session id")
	}
	if err := validateSeatIDs(req.SeatIDs); err != nil {
		return nil, err
	}

	booking, err := uc.repo.CreateBooking(ctx, userID, req)
	if err != nil {
		metrics.IncBookingErrors("create_booking_failed")
		return nil, err
	}
	metrics.IncActiveBookings()

	booking.ExpiresAt = booking.CreatedAt.Add(reservationHoldDuration)
	booking.SeatIDs = append([]int(nil), req.SeatIDs...)

	return booking, nil
}

func (uc *useCase) GetBookingByID(ctx context.Context, userID int, bookingID string) (*cinema.Booking, error) {
	if userID <= 0 {
		return nil, cinema.NewValidationError("user id")
	}
	if bookingID == "" {
		return nil, cinema.NewValidationError("booking id")
	}

	return uc.repo.GetBookingByID(ctx, userID, bookingID)
}

func (uc *useCase) RequestBookingPayment(ctx context.Context, userID int, bookingID string) (*cinema.Booking, error) {
	if userID <= 0 {
		metrics.IncBookingErrors("invalid_user_id")
		return nil, cinema.NewValidationError("user id")
	}
	if bookingID == "" {
		metrics.IncBookingErrors("empty_booking_id")
		return nil, cinema.NewValidationError("booking id")
	}
	if uc.publisher == nil {
		return nil, cinema.ErrInternal
	}

	booking, err := uc.repo.GetBookingByID(ctx, userID, bookingID)
	if err != nil {
		metrics.IncBookingErrors("purchase_booking_failed")
		return nil, err
	}
	if booking.Status != "pending" {
		return nil, cinema.NewConflictError("booking is not available for payment")
	}

	event := cinemaKafka.PaymentRequestedEvent{
		BookingID:   booking.ID,
		UserID:      userID,
		TotalPrice:   booking.TotalPrice,
		RequestedAt: time.Now().UTC(),
	}
	if err := uc.publisher.PublishPaymentRequested(ctx, event); err != nil {
		metrics.IncBookingErrors("payment_publish_failed")
		return nil, cinema.NewDatabaseError(fmt.Sprintf("publish payment request: %v", err))
	}

	booking.Status = "pending"
	return booking, nil
}

func (uc *useCase) CancelBooking(ctx context.Context, userID int, bookingID string) (*cinema.Booking, error) {
	if userID <= 0 {
		metrics.IncBookingErrors("invalid_user_id")
		return nil, cinema.NewValidationError("user id")
	}
	if bookingID == "" {
		metrics.IncBookingErrors("empty_booking_id")
		return nil, cinema.NewValidationError("booking id")
	}

	booking, err := uc.repo.GetBookingByID(ctx, userID, bookingID)
	if err != nil {
		metrics.IncBookingErrors("cancel_booking_failed")
		return nil, err
	}
	if booking.Status != "pending" {
		return nil, cinema.NewConflictError("only pending bookings can be cancelled")
	}

	if err := uc.repo.CancelBooking(ctx, userID, bookingID); err != nil {
		metrics.IncBookingErrors("cancel_booking_failed")
		return nil, err
	}
	metrics.DecActiveBookings()

	booking.Status = "cancelled"
	return booking, nil
}

func (uc *useCase) GetAllMyBookings(ctx context.Context, userID int) ([]cinema.Booking, error) {
	if userID <= 0 {
		return nil, cinema.NewValidationError("user id")
	}

	return uc.repo.GetAllMyBookings(ctx, userID)
}

func (uc *useCase) GetReservedSeatIDsBySession(ctx context.Context, sessionID int) ([]int, error) {
	if sessionID <= 0 {
		return nil, cinema.NewValidationError("session id")
	}
	return uc.repo.GetReservedSeatIDsBySession(ctx, sessionID)
}

// Auth методы
func (uc *useCase) Register(ctx context.Context, req *cinema.RegisterRequest) (*cinema.AuthResponse, error) {
	if req == nil || req.Email == "" || req.Password == "" {
		return nil, cinema.ErrInvalidInput
	}

	existing, _ := uc.repo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, cinema.ErrEmailAlreadyExists
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, cinema.NewDatabaseError("failed to hash password")
	}

	user := &cinema.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		IsActive:     true,
	}

	userID, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	if uc.tokenManager == nil {
		return nil, cinema.ErrInternal
	}

	token, err := uc.tokenManager.GenerateToken(userID, user.Email)
	if err != nil {
		return nil, cinema.NewDatabaseError("failed to generate token")
	}

	return &cinema.AuthResponse{
		Token:     token,
		UserID:    userID,
		Email:     req.Email,
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
	}, nil
}

func (uc *useCase) Login(ctx context.Context, req *cinema.LoginRequest) (*cinema.AuthResponse, error) {
	if req == nil || req.Email == "" || req.Password == "" {
		return nil, cinema.ErrInvalidInput
	}

	user, err := uc.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, cinema.NewNotFoundError("user")
	}

	if !security.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, cinema.AppError{
            Code:    "INVALID_CREDENTIALS",
            Message: "Invalid email or password",
            Status:  401,
        }
	}

	if uc.tokenManager == nil {
		return nil, cinema.ErrInternal
	}

	token, err := uc.tokenManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, cinema.NewDatabaseError("failed to generate token")
	}

	return &cinema.AuthResponse{
		Token:     token,
		UserID:    user.ID,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		}, nil
}

// Пользователи
func (uc *useCase) ListUsers(ctx context.Context) ([]cinema.User, error) {
	return uc.repo.GetAllUsers(ctx)
}

func (uc *useCase) GetMe(ctx context.Context, userID int) (*cinema.User, error) {
	if userID <= 0 {
		return nil, cinema.NewValidationError("user id")
	}
	return uc.repo.GetUserByID(ctx, userID)
}

// Билеты
func (uc *useCase) ListTickets(ctx context.Context, userID int) ([]cinema.Ticket, error) {
	if userID <= 0 {
		return nil, cinema.NewValidationError("user id")
	}
	return uc.repo.GetAllTickets(ctx, userID)
}

func (uc *useCase) GetTicketByID(ctx context.Context, userID int, ticketID int) (*cinema.Ticket, error) {
	if userID <= 0 {
		return nil, cinema.NewValidationError("user id")
	}
	if ticketID <= 0 {
		return nil, cinema.NewValidationError("ticket id")
	}
	return uc.repo.GetTicketByID(ctx, userID, ticketID)
}

// Тестовые методы
func (uc *useCase) GetTestMessage(ctx context.Context) (string, error) {
	return uc.repo.GetTestMessage(ctx)
}

func (uc *useCase) GetSlowMessage(ctx context.Context) (string, error) {
	return uc.repo.GetSlowMessage(ctx)
}
