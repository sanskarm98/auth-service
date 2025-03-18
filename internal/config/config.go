package config

import (
	"os"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Port            string
	JWTSecret       string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	// Default to 15 minutes for access token
	accessTokenExp := 15 * time.Minute

	// Default to 7 days for refresh token
	refreshTokenExp := 7 * 24 * time.Hour

	return &Config{
		Port:            port,
		JWTSecret:       jwtSecret,
		AccessTokenExp:  accessTokenExp,
		RefreshTokenExp: refreshTokenExp,
	}
}
