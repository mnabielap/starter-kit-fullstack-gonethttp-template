package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Code    int         `json:"code,omitempty"`
	Message interface{} `json:"message,omitempty"` // Can be string or map of errors
	Data    interface{} `json:"data,omitempty"`
}

// JSON writes a JSON response with a specific status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// Error writes a standardized error response
func Error(w http.ResponseWriter, status int, message interface{}) {
	JSON(w, status, APIResponse{
		Code:    status,
		Message: message,
	})
}

// Success writes a standardized success response (optional wrapper)
func Success(w http.ResponseWriter, status int, data interface{}) {
	// If the data is already a map/struct, just return it directly like the PHP controller
	JSON(w, status, data)
}