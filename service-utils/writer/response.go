package writer

import (
	"encoding/json"
	"net/http"
)

// Response defines the standard JSON structure returned by Gateway.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// JSON writes a JSON response with the given HTTP status code.
func JSON(w http.ResponseWriter, statusCode int, response Response) {
	// Tell the client that the response body contains JSON.
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code before writing the response body.
	w.WriteHeader(statusCode)

	// Encode the response structure as JSON.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// The response may already be partially written, so log the failure.
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// Success writes a successful HTTP response.
func Success(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error writes an error HTTP response.
func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, Response{
		Success: false,
		Message: message,
		Data:    nil,
	})
}
