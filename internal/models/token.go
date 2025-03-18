package models

import "github.com/golang-jwt/jwt/v5"

// TokenPair contains access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// RefreshRequest represents the request payload for refreshing tokens
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
