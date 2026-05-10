package usecase

import (
	"context"
	"time"

	"github.com/ElliAbby/go_cinema_system/internal/cinemaService"
	"github.com/ElliAbby/go_cinema_system/internal/jwt"
  "github.com/ElliAbby/go_cinema_system/internal/security"

)

type useCase struct {
	repo cinemaService.Repository
}

const reservationHoldDuration = 15 * time.Minute

func New(r cinemaService.Repository) *useCase {
	return &useCase{repo: r}
}

// Фильмы
func (uc *useCase) GetAllMovies(ctx context.Context) ([]cinemaService.Movie, error) {
	return uc.repo.GetAllMovies(ctx)
}

func (uc *useCase) GetMovieByID(ctx context.Context, id int) (*cinemaService.Movie, error) {
	if id <= 0 {
		return nil, cinemaService.NewValidationError("movie id")
	}
	return uc.repo.GetMovieByID(ctx, id)
}

func (uc *useCase) CreateMovie(ctx context.Context, req *cinemaService.CreateMovieRequest) (int, error) {
	if req == nil {
		return 0, cinemaService.ErrInvalidInput
	}
	if req.Title == "" {
		return 0, cinemaService.NewValidationError("title")
	}
	if req.Duration <= 0 {
		return 0, cinemaService.NewValidationError("duration")
	}

	movie := &cinemaService.Movie{
		Title:       req.Title,
		Duration:    req.Duration,
		Rating:      req.Rating,
		Description: req.Description,
	}
	return uc.repo.CreateMovie(ctx, movie)
}

func (uc *useCase) UpdateMovie(ctx context.Context, id int, req *cinemaService.CreateMovieRequest) error {
	if id <= 0 {
		return cinemaService.NewValidationError("movie id")
	}
	if req == nil {
		return cinemaService.ErrInvalidInput
	}
	if req.Title == "" {
		return cinemaService.NewValidationError("title")
	}
	if req.Duration <= 0 {
		return cinemaService.NewValidationError("duration")
	}

	movie := &cinemaService.Movie{
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
		return cinemaService.NewValidationError("movie id")
	}
	return uc.repo.DeleteMovie(ctx, id)
}

// Кинотеатры
func (uc *useCase) GetAllCinemas(ctx context.Context) ([]cinemaService.Cinema, error) {
	return uc.repo.GetAllCinemas(ctx)
}

func (uc *useCase) GetCinemaByID(ctx context.Context, id int) (*cinemaService.Cinema, error) {
	if id <= 0 {
		return nil, cinemaService.NewValidationError("cinema id")
	}
	return uc.repo.GetCinemaByID(ctx, id)
}

func (uc *useCase) CreateCinema(ctx context.Context, req *cinemaService.CreateCinemaRequest) (int, error) {
	if req == nil {
		return 0, cinemaService.ErrInvalidInput
	}
	if req.Name == "" {
		return 0, cinemaService.NewValidationError("name")
	}
	if req.Address == "" {
		return 0, cinemaService.NewValidationError("address")
	}

	cinema := &cinemaService.Cinema{
		Name:    req.Name,
		Address: req.Address,
	}
	return uc.repo.CreateCinema(ctx, cinema)
}

// Залы
func (uc *useCase) GetHallsByCinema(ctx context.Context, cinemaID int) ([]cinemaService.Hall, error) {
	if cinemaID <= 0 {
		return nil, cinemaService.NewValidationError("cinema id")
	}
	return uc.repo.GetHallsByCinema(ctx, cinemaID)
}

// Сессии
func (uc *useCase) GetAllSessions(ctx context.Context) ([]cinemaService.Session, error) {
	return uc.repo.GetAllSessions(ctx)
}

func (uc *useCase) GetSessionByID(ctx context.Context, id int) (*cinemaService.Session, error) {
	if id <= 0 {
		return nil, cinemaService.NewValidationError("session id")
	}
	return uc.repo.GetSessionByID(ctx, id)
}

func (uc *useCase) GetSessionsByMovie(ctx context.Context, movieID int) ([]cinemaService.Session, error) {
	if movieID <= 0 {
		return nil, cinemaService.NewValidationError("movie id")
	}
	return uc.repo.GetSessionsByMovie(ctx, movieID)
}

func (uc *useCase) CreateSession(ctx context.Context, req *cinemaService.CreateSessionRequest) (int, error) {
	if req == nil {
		return 0, cinemaService.ErrInvalidInput
	}
	if req.MovieID <= 0 {
		return 0, cinemaService.NewValidationError("movie id")
	}
	if req.HallID <= 0 {
		return 0, cinemaService.NewValidationError("hall id")
	}
	if req.PriceBase < 0 {
		return 0, cinemaService.NewValidationError("price base")
	}

	session := &cinemaService.Session{
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
		return cinemaService.NewValidationError("seat ids")
	}

	seen := make(map[int]struct{}, len(seatIDs))
	for _, seatID := range seatIDs {
		if seatID <= 0 {
			return cinemaService.NewValidationError("seat id")
		}
		if _, exists := seen[seatID]; exists {
			return cinemaService.NewValidationError("seat ids")
		}
		seen[seatID] = struct{}{}
	}

	return nil
}

func (uc *useCase) CreateBooking(ctx context.Context, userID int, req *cinemaService.CreateBookingRequest) (*cinemaService.Booking, error) {
	if userID <= 0 {
		return nil, cinemaService.NewValidationError("user id")
	}
	if req == nil {
		return nil, cinemaService.ErrInvalidInput
	}
	if req.SessionID <= 0 {
		return nil, cinemaService.NewValidationError("session id")
	}
	if err := validateSeatIDs(req.SeatIDs); err != nil {
		return nil, err
	}

	booking, err := uc.repo.CreateBooking(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	booking.ExpiresAt = booking.CreatedAt.Add(reservationHoldDuration)
	booking.SeatIDs = append([]int(nil), req.SeatIDs...)
	return booking, nil
}

func (uc *useCase) PurchaseBooking(ctx context.Context, userID int, bookingID string) (*cinemaService.Booking, []cinemaService.Ticket, error) {
	if userID <= 0 {
		return nil, nil, cinemaService.NewValidationError("user id")
	}
	if bookingID == "" {
		return nil, nil, cinemaService.NewValidationError("booking id")
	}

	booking, tickets, err := uc.repo.PurchaseBooking(ctx, userID, bookingID)
	if err != nil {
		return nil, nil, err
	}

	if booking != nil {
		booking.SeatIDs = make([]int, 0, len(tickets))
		for _, ticket := range tickets {
			booking.SeatIDs = append(booking.SeatIDs, ticket.SeatID)
		}
	}

	return booking, tickets, nil
}

func (uc *useCase) GetAllMyBookings(ctx context.Context, userID int) ([]cinemaService.Booking, error) {
	if userID <= 0 {
		return nil, cinemaService.NewValidationError("user id")
	}

	return uc.repo.GetAllMyBookings(ctx, userID)
}

// Auth методы
func (uc *useCase) Register(ctx context.Context, req *cinemaService.RegisterRequest) (*cinemaService.AuthResponse, error) {
	if req == nil || req.Email == "" || req.Password == "" {
		return nil, cinemaService.ErrInvalidInput
	}

	existing, _ := uc.repo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, cinemaService.ErrEmailAlreadyExists
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, cinemaService.NewDatabaseError("failed to hash password")
	}

	user := &cinemaService.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		IsActive:     true,
	}

	userID, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(userID, user.Email)
	if err != nil {
		return nil, cinemaService.NewDatabaseError("failed to generate token")
	}

	return &cinemaService.AuthResponse{
		Token:     token,
		UserID:    userID,
		Email:     req.Email,
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
	}, nil
}

func (uc *useCase) Login(ctx context.Context, req *cinemaService.LoginRequest) (*cinemaService.AuthResponse, error) {
	if req == nil || req.Email == "" || req.Password == "" {
		return nil, cinemaService.ErrInvalidInput
	}

	user, err := uc.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, cinemaService.NewNotFoundError("user")
	}

	if !security.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, cinemaService.AppError{
            Code:    "INVALID_CREDENTIALS",
            Message: "Invalid email or password",
            Status:  401,
        }
	}

	token, err := jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, cinemaService.NewDatabaseError("failed to generate token")
	}

	return &cinemaService.AuthResponse{
		Token:     token,
		UserID:    user.ID,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		}, nil
}

// Пользователи
func (uc *useCase) ListUsers(ctx context.Context) ([]cinemaService.User, error) {
	return uc.repo.GetAllUsers(ctx)
}

func (uc *useCase) GetMe(ctx context.Context, userID int) (*cinemaService.User, error) {
	if userID <= 0 {
		return nil, cinemaService.NewValidationError("user id")
	}
	return uc.repo.GetUserByID(ctx, userID)
}

// Билеты
func (uc *useCase) ListTickets(ctx context.Context, userID int) ([]cinemaService.Ticket, error) {
	if userID <= 0 {
		return nil, cinemaService.NewValidationError("user id")
	}
	return uc.repo.GetAllTickets(ctx, userID)
}

func (uc *useCase) GetTicketByID(ctx context.Context, userID int, ticketID int) (*cinemaService.Ticket, error) {
	if userID <= 0 {
		return nil, cinemaService.NewValidationError("user id")
	}
	if ticketID <= 0 {
		return nil, cinemaService.NewValidationError("ticket id")
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
