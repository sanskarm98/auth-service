package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sanskarm98/auth-service/internal/models"
)

// SendJSONResponse sends a JSON response with the provided status code
func SendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log the error or return an HTTP error response
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

// SendErrorResponse sends a standardized error response
func SendErrorResponse(w http.ResponseWriter, status int, message string) {
	errResponse := models.ErrorResponse{
		Status:  status,
		Message: message,
	}
	SendJSONResponse(w, status, errResponse)
}

// ExtractTokenFromHeader extracts JWT from Authorization header
func ExtractTokenFromHeader(r *http.Request) string {
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
