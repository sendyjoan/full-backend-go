package errors

import (
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

// Common error types
var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrEmailNotFound       = errors.New("email not found")
	ErrInvalidOTP          = errors.New("invalid or expired otp")
	ErrUserNotFound        = errors.New("user not found")
	ErrTokenExpired        = errors.New("token expired")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrValidationFailed    = errors.New("validation failed")
	ErrInternalServer      = errors.New("internal server error")
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	CodeInvalidCredentials  ErrorCode = "INVALID_CREDENTIALS"
	CodeInvalidRefreshToken ErrorCode = "INVALID_REFRESH_TOKEN"
	CodeEmailNotFound       ErrorCode = "EMAIL_NOT_FOUND"
	CodeInvalidOTP          ErrorCode = "INVALID_OTP"
	CodeUserNotFound        ErrorCode = "USER_NOT_FOUND"
	CodeTokenExpired        ErrorCode = "TOKEN_EXPIRED"
	CodeUnauthorized        ErrorCode = "UNAUTHORIZED"
	CodeValidationFailed    ErrorCode = "VALIDATION_FAILED"
	CodeInternalServer      ErrorCode = "INTERNAL_SERVER_ERROR"
)

// AppError represents application-specific errors
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// New creates a new AppError
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Newf creates a new AppError with formatted message
func Newf(code ErrorCode, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// WithDetails adds details to an AppError
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// ToHumaError converts AppError to Huma error
func (e *AppError) ToHumaError() error {
	switch e.Code {
	case CodeInvalidCredentials, CodeInvalidRefreshToken, CodeUnauthorized:
		return huma.Error401Unauthorized(e.Message)
	case CodeEmailNotFound, CodeUserNotFound:
		return huma.Error404NotFound(e.Message)
	case CodeValidationFailed, CodeInvalidOTP:
		return huma.Error400BadRequest(e.Message)
	default:
		return huma.Error500InternalServerError(e.Message)
	}
}

// Helper functions for common errors
func InvalidCredentials() *AppError {
	return New(CodeInvalidCredentials, "Invalid username/email or password")
}

func InvalidRefreshToken() *AppError {
	return New(CodeInvalidRefreshToken, "Invalid or expired refresh token")
}

func EmailNotFound() *AppError {
	return New(CodeEmailNotFound, "Email address not found")
}

func InvalidOTP() *AppError {
	return New(CodeInvalidOTP, "Invalid or expired OTP code")
}

func UserNotFound() *AppError {
	return New(CodeUserNotFound, "User not found")
}

func TokenExpired() *AppError {
	return New(CodeTokenExpired, "Token has expired")
}

func Unauthorized() *AppError {
	return New(CodeUnauthorized, "Unauthorized access")
}

func ValidationFailed(details string) *AppError {
	return New(CodeValidationFailed, "Validation failed").WithDetails(details)
}

func InternalServer(details string) *AppError {
	return New(CodeInternalServer, "Internal server error").WithDetails(details)
}

// IsAppError checks if error is an AppError
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
