package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sanskarm98/auth-service/internal/models"
	"github.com/sanskarm98/auth-service/internal/store"
	"github.com/sanskarm98/auth-service/pkg/utils"
)

// Custom context key type to avoid collisions
type contextKey string

const (
	// ClaimsContextKey is the key for JWT claims in the request context
	ClaimsContextKey contextKey = "claims"
)

// AuthMiddleware handles JWT authentication for protected routes
type AuthMiddleware struct {
	authService AuthService
	tokenStore  store.TokenStore
}

// NewAuthMiddleware creates a new instance of AuthMiddleware
func NewAuthMiddleware(authService AuthService, tokenStore store.TokenStore) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		tokenStore:  tokenStore,
	}
}

// Authenticate is a middleware that verifies the access token in the Authorization header
func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			utils.SendErrorResponse(w, http.StatusUnauthorized, models.ErrTokenRequired)
			return
		}

		// Check if token is revoked
		if m.tokenStore.IsTokenRevoked(tokenString) {
			utils.SendErrorResponse(w, http.StatusUnauthorized, models.ErrTokenRevoked)
			return
		}

		// Parse and validate token
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				utils.SendErrorResponse(w, http.StatusUnauthorized, models.ErrTokenExpired)
				return
			}
			utils.SendErrorResponse(w, http.StatusUnauthorized, models.ErrInvalidToken)
			return
		}

		// Set claims in context and proceed
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// extractTokenFromHeader extracts JWT from Authorization header
func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}

	// Check if it's a bearer token
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// GetClaimsFromContext extracts claims from request context
func GetClaimsFromContext(ctx context.Context) (*models.Claims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(*models.Claims)
	return claims, ok
}
