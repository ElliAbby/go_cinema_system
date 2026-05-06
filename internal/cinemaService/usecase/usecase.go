package usecase

import (
	"context"

	"github.com/ElliAbby/go_cinema_system/internal/cinemaService"

)

type useCase struct {
	repo cinemaService.Repository
}

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

// Тестовые методы
func (uc *useCase) GetTestMessage(ctx context.Context) (string, error) {
	return uc.repo.GetTestMessage(ctx)
}

func (uc *useCase) GetSlowMessage(ctx context.Context) (string, error) {
	return uc.repo.GetSlowMessage(ctx)
}
