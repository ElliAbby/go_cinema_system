package cinemaService

import "fmt"

// AppError представляет ошибку приложения
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e AppError) Error() string {
	return e.Message
}

// Предопределенные ошибки
var (
	ErrNotFound = AppError{
		Code:    "NOT_FOUND",
		Message: "Resource not found",
		Status:  404,
	}
	ErrInvalidInput = AppError{
		Code:    "INVALID_INPUT",
		Message: "Invalid input data",
		Status:  400,
	}
	ErrConflict = AppError{
		Code:    "CONFLICT",
		Message: "Resource already exists",
		Status:  409,
	}
	ErrInternal = AppError{
		Code:    "INTERNAL_ERROR",
		Message: "Internal server error",
		Status:  500,
	}
)

// NewNotFoundError возвращает ошибку NotFound с кастомным сообщением
func NewNotFoundError(resource string) AppError {
	return AppError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s not found", resource),
		Status:  404,
	}
}

// NewValidationError возвращает ошибку валидации
func NewValidationError(field string) AppError {
	return AppError{
		Code:    "INVALID_INPUT",
		Message: fmt.Sprintf("Invalid %s", field),
		Status:  400,
	}
}

// NewDatabaseError возвращает ошибку БД
func NewDatabaseError(msg string) AppError {
	return AppError{
		Code:    "DATABASE_ERROR",
		Message: fmt.Sprintf("Database error: %s", msg),
		Status:  500,
	}
}
