package shared

import (
	"github.com/gofiber/fiber/v2"
)

// JoyStickError represents a custom error type for the JoyStick application.
type JoyStickError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface for JoyStickError.
func (e *JoyStickError) Error() string {
	return e.Message
}

// ErrorCode returns the error code for JoyStickError.
func (e *JoyStickError) ErrorCode() int {
	return e.Code
}

var (
	// ErrInvalidPasswordFormat is returned when the password does not meet the required format.
	ErrInvalidPasswordFormat = &JoyStickError{
		Code:    fiber.StatusBadRequest,
		Message: "password must be at least 8 characters long and contain at least one number, one uppercase letter, one lowercase letter, and one special character",
	}

	// ErrUsernameTooShort is returned when the username is too short.
	ErrUsernameTooShort = &JoyStickError{
		Code:    fiber.StatusBadRequest,
		Message: "username must be at least 3 characters long",
	}

	// ErrUserAlreadyExists is returned when a user with the given username already exists.
	ErrUserAlreadyExists = &JoyStickError{
		Code:    fiber.StatusConflict,
		Message: "user with this username already exists",
	}

	// ErrUserNotFound is returned when a user with the given username is not found.
	ErrUserNotFound = &JoyStickError{
		Code:    fiber.StatusNotFound,
		Message: "user not found",
	}

	// ErrInvalidCredentials is returned when the provided credentials are invalid.
	ErrInvalidCredentials = &JoyStickError{
		Code:    fiber.StatusUnauthorized,
		Message: "invalid username or password",
	}

	// ErrInternalServerError is returned for unexpected server errors.
	ErrInternalServerError = &JoyStickError{
		Code:    fiber.StatusInternalServerError,
		Message: "internal server error",
	}

	// ErrServiceNotFound is returned when a service with the given ID is not found.
	ErrServiceNotFound = &JoyStickError{
		Code:    fiber.StatusNotFound,
		Message: "service not found",
	}

	// ErrInvalidToken is returned when the provided token is invalid.
	ErrInvalidToken = &JoyStickError{
		Code:    fiber.StatusUnauthorized,
		Message: "invalid token",
	}

	// ErrServiceDeploymentNotFound is returned when a service deployment with the given ID is not found.
	ErrServiceDeploymentNotFound = &JoyStickError{
		Code:    fiber.StatusNotFound,
		Message: "service deployment not found",
	}
)
