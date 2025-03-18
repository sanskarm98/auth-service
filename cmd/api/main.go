package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sanskarm98/auth-service/internal/auth"
	"github.com/sanskarm98/auth-service/internal/config"
	"github.com/sanskarm98/auth-service/internal/handlers"
	"github.com/sanskarm98/auth-service/internal/store"
)

func main() {
	// Initialize configuration
	cfg := config.LoadConfig()

	// Initialize stores
	userStore := store.NewInMemoryUserStore()
	tokenStore := store.NewInMemoryTokenStore()

	// Initialize auth service
	authService := auth.NewJWTAuthService(cfg.JWTSecret, cfg.AccessTokenExp, cfg.RefreshTokenExp, tokenStore)

	// Initialize middleware
	authMiddleware := auth.NewAuthMiddleware(authService, tokenStore)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userStore, authService, tokenStore)
	userHandler := handlers.NewUserHandler(userStore)

	// Setup routes
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/api/auth/signup", authHandler.SignUp)
	mux.HandleFunc("/api/auth/signin", authHandler.SignIn)
	mux.HandleFunc("/api/auth/refresh", authHandler.RefreshToken)
	mux.HandleFunc("/api/auth/revoke", authMiddleware.Authenticate(authHandler.RevokeToken))
	mux.HandleFunc("/api/auth/verify", authMiddleware.Authenticate(authHandler.VerifyToken))

	// User routes
	mux.HandleFunc("/api/auth/me", authMiddleware.Authenticate(userHandler.GetUserInfo))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	fmt.Printf("Auth service starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
