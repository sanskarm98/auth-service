package models

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ErrorMessages holds constant error messages to be used across the application
const (
	ErrInvalidRequest      = "Invalid request payload"
	ErrInvalidCredentials  = "Invalid credentials"
	ErrEmailAlreadyExists  = "Email already registered"
	ErrInvalidToken        = "Invalid token"
	ErrTokenExpired        = "Token has expired"
	ErrTokenRevoked        = "Token has been revoked"
	ErrUserNotFound        = "User not found"
	ErrMethodNotAllowed    = "Method not allowed"
	ErrTokenRequired       = "Authorization token required"
	ErrInternalServerError = "Internal server error"
	ErrInvalidRefreshToken = "Invalid refresh token"
	ErrRequiredFields      = "Required fields missing"
)
