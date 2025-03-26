package response

import (
	"encoding/json"
	"net/http"
)

type APIResponseSuccess struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
}

type APIResponseError struct {
	Status     string        `json:"status"`
	StatusCode int           `json:"status_code"`
	Message    string        `json:"message,omitempty"`
	Errors     []ErrorDetail `json:"errors,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// writeResponse is a generic utility for writing responses
func writeResponseSuccess(w http.ResponseWriter, status string, message string, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponseSuccess{
		Status:     status,
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func writeResponseError(w http.ResponseWriter, status string, message string, errors []ErrorDetail, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponseError{
		Status:     status,
		StatusCode: statusCode,
		Message:    message,
		Errors:     errors,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func sendError(w http.ResponseWriter, message string, errors []ErrorDetail, statusCode int) {
	writeResponseError(w, "error", message, errors, statusCode)
}

func sendSuccess(w http.ResponseWriter, message string, data any, statusCode int) {
	writeResponseSuccess(w, "success", message, data, statusCode)
}

// Informational responses
func Send100Continue(w http.ResponseWriter, message string) {
	writeResponseSuccess(w, "info", message, nil, http.StatusContinue)
}

func Send101SwitchingProtocols(w http.ResponseWriter, message string) {
	writeResponseSuccess(w, "info", message, nil, http.StatusSwitchingProtocols)
}

func Send102Processing(w http.ResponseWriter, message string) {
	writeResponseSuccess(w, "info", message, nil, http.StatusProcessing)
}

func Send103EarlyHints(w http.ResponseWriter, message string) {
	writeResponseSuccess(w, "info", message, nil, http.StatusEarlyHints)
}

// Success responses
func Send200OK(w http.ResponseWriter, message string, data any) {
	sendSuccess(w, message, data, http.StatusOK)
}

func Send201Created(w http.ResponseWriter, message string, data any) {
	sendSuccess(w, message, data, http.StatusCreated)
}

func Send202Accepted(w http.ResponseWriter, message string) {
	sendSuccess(w, message, "", http.StatusAccepted)
}

func Send203NonAuthoritativeInfo(w http.ResponseWriter, message string) {
	sendSuccess(w, message, "", http.StatusNonAuthoritativeInfo)
}

func Send204NoContent(w http.ResponseWriter) {
	writeResponseSuccess(w, "success", "No Content", nil, http.StatusNoContent)
}

func Send205ResetContent(w http.ResponseWriter, message string) {
	writeResponseSuccess(w, "success", message, nil, http.StatusResetContent)
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

func SendCustomResponse(w http.ResponseWriter, errors []ErrorDetail, statusCode int) {
	sendError(w, "err", errors, statusCode)
}
