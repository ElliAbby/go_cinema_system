package cinemaService

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

	// Места
	GetSeatsByHall(ctx context.Context, hallID int) ([]Seat, error)
	GetSeatByID(ctx context.Context, id int) (*Seat, error)

	// Тестовые методы
	GetTestMessage(ctx context.Context) (string, error)
	GetSlowMessage(ctx context.Context) (string, error)
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

	// Сессии
	GetAllSessions(ctx context.Context) ([]Session, error)
	GetSessionByID(ctx context.Context, id int) (*Session, error)
	GetSessionsByMovie(ctx context.Context, movieID int) ([]Session, error)
	CreateSession(ctx context.Context, req *CreateSessionRequest) (int, error)

	// Тестовые методы
	GetTestMessage(ctx context.Context) (string, error)
	GetSlowMessage(ctx context.Context) (string, error)
}
