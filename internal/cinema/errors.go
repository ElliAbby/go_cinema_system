package cinema

import (
	apperror "github.com/ElliAbby/go_cinema_system/internal/platform/apperror"
)

type AppError = apperror.AppError

var (
	ErrNotFound           = apperror.ErrNotFound
	ErrInvalidInput       = apperror.ErrInvalidInput
	ErrConflict           = apperror.ErrConflict
	ErrInternal           = apperror.ErrInternal
	ErrEmailAlreadyExists = apperror.ErrEmailAlreadyExists
)

func NewNotFoundError(resource string) AppError {
	return apperror.NewNotFoundError(resource)
}

func NewValidationError(field string) AppError {
	return apperror.NewValidationError(field)
}

func NewConflictError(message string) AppError {
	return apperror.NewConflictError(message)
}

func NewDatabaseError(msg string) AppError {
	return apperror.NewDatabaseError(msg)
}
