package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse represents a standardized API error response format
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Status  int        `json:"status"`

}

// SendErrorResponse writes a JSON error response with proper headers
func SendErrorResponse(w http.ResponseWriter, message string, status int) error {
	response := ErrorResponse{
		Success: false,
		Message: message,
		Status:  status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("failed to encode error response: %w", err)
	}
	return nil
}