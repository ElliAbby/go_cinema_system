package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

)

type repo struct {
	postgresDB *sqlx.DB
}

func New(postgresDB *sqlx.DB) *repo {
	return &repo{postgresDB: postgresDB}
}

func (r *repo) GetTestMessage(ctx context.Context) (string, error) {
	return "Hello!", nil
}

func (r *repo) GetSlowMessage(ctx context.Context) (string, error) {
	return "Slooooow text", nil
}
