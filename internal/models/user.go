package models

import "time"

// User represents the user model
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Don't return password in responses
	CreatedAt time.Time `json:"created_at"`
}

// SignupRequest represents the request payload for user registration
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SigninRequest represents the request payload for user authentication
type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// NewUserResponse creates a new UserResponse from a User model
func NewUserResponse(user User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
