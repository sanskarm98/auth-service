package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sanskarm98/auth-service/internal/auth"
	"github.com/sanskarm98/auth-service/internal/models"
	"github.com/sanskarm98/auth-service/internal/store"
	"github.com/sanskarm98/auth-service/pkg/utils"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	userStore   store.UserStore
	authService auth.AuthService
	tokenStore  store.TokenStore
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(userStore store.UserStore, authService auth.AuthService, tokenStore store.TokenStore) *AuthHandler {
	return &AuthHandler{
		userStore:   userStore,
		authService: authService,
		tokenStore:  tokenStore,
	}
}

// SignUp handles user registration
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrMethodNotAllowed)
		return
	}

	// Parse request
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrInvalidRequest)
		return
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrRequiredFields)
		return
	}

	// Create user
	user, err := h.userStore.Create(req.Email, req.Password)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusConflict, err.Error())
		return
	}

	// Return user data
	utils.SendJSONResponse(w, http.StatusCreated, models.NewUserResponse(user))
}

// SignIn handles user authentication
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrMethodNotAllowed)
		return
	}

	// Parse request
	var req models.SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrInvalidRequest)
		return
	}

	// Authenticate user
	user, authenticated := h.userStore.Authenticate(req.Email, req.Password)
	if !authenticated {
		utils.SendErrorResponse(w, http.StatusUnauthorized, models.ErrInvalidCredentials)
		return
	}

	// Generate token pair
	tokenPair, err := h.authService.GenerateTokenPair(user)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrInternalServerError)
		return
	}

	// Return tokens
	utils.SendJSONResponse(w, http.StatusOK, tokenPair)
}

// RefreshToken issues new tokens based on a valid refresh token
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrMethodNotAllowed)
		return
	}

	// Parse request
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrInvalidRequest)
		return
	}

	// Validate refresh token
	userID, exists := h.tokenStore.GetUserIDByRefreshToken(req.RefreshToken)
	if !exists {
		utils.SendErrorResponse(w, http.StatusUnauthorized, models.ErrInvalidRefreshToken)
		return
	}

	// Get user
	user, exists := h.userStore.GetByID(userID)
	if !exists {
		utils.SendErrorResponse(w, http.StatusUnauthorized, models.ErrUserNotFound)
		return
	}

	// Delete old refresh token
	h.tokenStore.DeleteRefreshToken(req.RefreshToken)

	// Generate new token pair
	tokenPair, err := h.authService.GenerateTokenPair(user)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrInternalServerError)
		return
	}

	// Return tokens
	utils.SendJSONResponse(w, http.StatusOK, tokenPair)
}

// RevokeToken revokes an access token
func (h *AuthHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrMethodNotAllowed)
		return
	}

	// Extract token from header
	tokenString := utils.ExtractTokenFromHeader(r)
	if tokenString == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrTokenRequired)
		return
	}

	// Revoke token
	h.tokenStore.RevokeToken(tokenString)

	// Return success
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Token revoked successfully"})
}

// VerifyToken simply confirms that a token is valid
func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	// Only for demonstration - token verification is done by middleware
	claims, _ := auth.GetClaimsFromContext(r.Context())

	response := map[string]interface{}{
		"message": "Token verified successfully",
		"user_id": claims.UserID,
		"email":   claims.Email,
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}
