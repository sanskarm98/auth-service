package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sanskarm98/auth-service/internal/models"
	"github.com/sanskarm98/auth-service/internal/store"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	GenerateTokenPair(user models.User) (models.TokenPair, error)
	ValidateToken(tokenString string) (*models.Claims, error)
}

// JWTAuthService implements AuthService with JWT tokens
type JWTAuthService struct {
	jwtSecret       string
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
	tokenStore      store.TokenStore
}

// NewJWTAuthService creates a new instance of JWTAuthService
func NewJWTAuthService(
	jwtSecret string,
	accessTokenExp time.Duration,
	refreshTokenExp time.Duration,
	tokenStore store.TokenStore,
) *JWTAuthService {
	return &JWTAuthService{
		jwtSecret:       jwtSecret,
		accessTokenExp:  accessTokenExp,
		refreshTokenExp: refreshTokenExp,
		tokenStore:      tokenStore,
	}
}

// GenerateTokenPair creates a new access and refresh token pair
func (s *JWTAuthService) GenerateTokenPair(user models.User) (models.TokenPair, error) {
	// Create access token
	accessExp := time.Now().Add(s.accessTokenExp)
	accessClaims := models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return models.TokenPair{}, err
	}

	// Create refresh token (simple UUID)
	refreshTokenString := uuid.New().String()

	// Store refresh token
	s.tokenStore.StoreRefreshToken(refreshTokenString, user.ID)

	return models.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

// ValidateToken validates a JWT token and returns its claims
func (s *JWTAuthService) ValidateToken(tokenString string) (*models.Claims, error) {
	// Parse and validate token
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	// Handle parsing errors
	if err != nil {
		return nil, err
	}

	// Validate token
	if !token.Valid {
		return nil, fmt.Errorf(models.ErrInvalidToken)
	}

	return claims, nil
}
