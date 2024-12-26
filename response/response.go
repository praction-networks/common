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

// writeResponse is a generic utility for writing responses
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
	writeResponse(w, "error", message, nil, errors, statusCode)
}

func sendSuccess(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	writeResponse(w, "success", message, data, nil, statusCode)
}

// Informational responses
func Send100Continue(w http.ResponseWriter, message string) {
	writeResponse(w, "info", message, nil, nil, http.StatusContinue)
}

func Send101SwitchingProtocols(w http.ResponseWriter, message string) {
	writeResponse(w, "info", message, nil, nil, http.StatusSwitchingProtocols)
}

func Send102Processing(w http.ResponseWriter, message string) {
	writeResponse(w, "info", message, nil, nil, http.StatusProcessing)
}

func Send103EarlyHints(w http.ResponseWriter, message string) {
	writeResponse(w, "info", message, nil, nil, http.StatusEarlyHints)
}

// Success responses
func Send200OK(w http.ResponseWriter, message string, data interface{}) {
	sendSuccess(w, message, data, http.StatusOK)
}

func Send201Created(w http.ResponseWriter, message string, data interface{}) {
	sendSuccess(w, message, data, http.StatusCreated)
}

func Send202Accepted(w http.ResponseWriter, message string, data interface{}) {
	sendSuccess(w, message, data, http.StatusAccepted)
}

func Send203NonAuthoritativeInfo(w http.ResponseWriter, message string, data interface{}) {
	sendSuccess(w, message, data, http.StatusNonAuthoritativeInfo)
}

func Send204NoContent(w http.ResponseWriter) {
	writeResponse(w, "success", "No Content", nil, nil, http.StatusNoContent)
}

func Send205ResetContent(w http.ResponseWriter, message string) {
	writeResponse(w, "success", message, nil, nil, http.StatusResetContent)
}

func Send206PartialContent(w http.ResponseWriter, message string, data interface{}) {
	sendSuccess(w, message, data, http.StatusPartialContent)
}

func Send207MultiStatus(w http.ResponseWriter, message string, data interface{}) {
	sendSuccess(w, message, data, http.StatusMultiStatus)
}

func Send208AlreadyReported(w http.ResponseWriter, message string, data interface{}) {
	sendSuccess(w, message, data, http.StatusAlreadyReported)
}

// Redirection responses
// Add similar handlers here for 300-series responses if needed

// Client error responses
func Send400BadRequest(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "request", Message: message}}, http.StatusBadRequest)
}

func Send401Unauthorized(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "authentication", Message: message}}, http.StatusUnauthorized)
}

func Send403Forbidden(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "authorization", Message: message}}, http.StatusForbidden)
}

func Send404NotFound(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "resource", Message: message}}, http.StatusNotFound)
}

func Send410Gone(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "resource", Message: message}}, http.StatusGone)
}

func Send409Conflict(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "conflict", Message: message}}, http.StatusConflict)
}

func Send415UnsupportedMediaType(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "media_type", Message: message}}, http.StatusUnsupportedMediaType)
}

// Server error responses
func Send500InternalServerError(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "server", Message: message}}, http.StatusInternalServerError)
}

func Send503ServiceUnavailable(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "server", Message: message}}, http.StatusServiceUnavailable)
}

func SendCustomError(w http.ResponseWriter, message string, errors []ErrorDetail, statusCode int) {
	sendError(w, message, errors, statusCode)
}

// calculateCount calculates the count for the data field
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
