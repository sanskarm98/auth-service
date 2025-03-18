package handlers

import (
	"net/http"

	"github.com/sanskarm98/auth-service/internal/auth"
	"github.com/sanskarm98/auth-service/internal/models"
	"github.com/sanskarm98/auth-service/internal/store"
	"github.com/sanskarm98/auth-service/pkg/utils"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userStore store.UserStore
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userStore store.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

// GetUserInfo returns the current user's details
func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, ok := auth.GetClaimsFromContext(r.Context())
	if !ok {
		utils.SendErrorResponse(w, http.StatusUnauthorized, models.ErrInvalidToken)
		return
	}

	// Get user
	user, exists := h.userStore.GetByID(claims.UserID)
	if !exists {
		utils.SendErrorResponse(w, http.StatusNotFound, models.ErrUserNotFound)
		return
	}

	// Return user data
	utils.SendJSONResponse(w, http.StatusOK, models.NewUserResponse(user))
}
