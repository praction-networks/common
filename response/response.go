package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Count   *int          `json:"count,omitempty"`
	Data    interface{}   `json:"data,omitempty"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// WriteResponse is a generic utility for writing responses
func writeResponse(w http.ResponseWriter, status string, message string, data interface{}, errors []ErrorDetail, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	count := calculateCount(data)

	response := APIResponse{
		Status:  status,
		Message: message,
		Count:   count,
		Data:    data,
		Errors:  errors,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func sendError(w http.ResponseWriter, message string, errors []ErrorDetail, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Status:  "error",
		Message: message, // Optional general error message
		Errors:  errors,  // Detailed error list
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"status":"error","message":"Failed to encode error response"}`, http.StatusInternalServerError)
	}
}

func sendSuccess(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	writeResponse(w, "success", message, data, nil, statusCode)
}

func SendCreated(w http.ResponseWriter, message string, data interface{}) {
	sendSuccess(w, message, data, http.StatusCreated)
}

func SendBadRequest(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "request", Message: message}}, http.StatusBadRequest)
}

func SendUnauthorized(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "authentication", Message: message}}, http.StatusUnauthorized)
}

func SendForbidden(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "authorization", Message: message}}, http.StatusForbidden)
}

func SendNotFound(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "resource", Message: message}}, http.StatusNotFound)
}

func SendGone(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "resource", Message: message}}, http.StatusGone)
}

func SendConflict(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "conflict", Message: message}}, http.StatusConflict)
}

func SendUnsupportedMediaType(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "media_type", Message: message}}, http.StatusUnsupportedMediaType)
}

func SendInternalServerError(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "server", Message: message}}, http.StatusInternalServerError)
}

func SendServiceUnavailable(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "server", Message: message}}, http.StatusServiceUnavailable)
}

func SendCustomError(w http.ResponseWriter, message string, errors []ErrorDetail, statusCode int) {
	sendError(w, message, errors, statusCode)
}

func calculateCount(data interface{}) *int {
	switch v := data.(type) {
	case []interface{}:
		count := len(v)
		return &count
	case []map[string]interface{}:
		count := len(v)
		return &count
	default:
		return nil
	}
}
