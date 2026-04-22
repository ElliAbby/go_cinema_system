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

func (uc *useCase) GetTestMessage(ctx context.Context) (string, error) {
	return uc.repo.GetTestMessage(ctx)
}

func (uc *useCase) GetSlowMessage(ctx context.Context) (string, error) {
	return uc.repo.GetSlowMessage(ctx)
}
