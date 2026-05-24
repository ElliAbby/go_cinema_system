package apperror

import "fmt"

// AppError represents an application-level error suitable for API responses.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e AppError) Error() string {
	return e.Message
}

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
	ErrEmailAlreadyExists = AppError{
		Code:    "USER_EXISTS",
		Message: "User with this email already exists",
		Status:  409,
	}
)

func NewNotFoundError(resource string) AppError {
	return AppError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s not found", resource),
		Status:  404,
	}
}

func NewValidationError(field string) AppError {
	return AppError{
		Code:    "INVALID_INPUT",
		Message: fmt.Sprintf("Invalid %s", field),
		Status:  400,
	}
}

func NewConflictError(message string) AppError {
	return AppError{
		Code:    "CONFLICT",
		Message: message,
		Status:  409,
	}
}

func NewDatabaseError(msg string) AppError {
	return AppError{
		Code:    "DATABASE_ERROR",
		Message: "Database operation failed",
		Status:  500,
	}
}