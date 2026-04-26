package cinemaService

import "context"


type UseCase interface {
	GetTestMessage(ctx context.Context) (string, error)
	GetSlowMessage(ctx context.Context) (string, error)
}

type Repository interface {
	GetTestMessage(ctx context.Context) (string, error)
	GetSlowMessage(ctx context.Context) (string, error)
}
