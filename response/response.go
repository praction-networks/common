package response

import (
	"encoding/json"
	"net/http"
)

type APIResponseSuccess struct {
	Status      string     `json:"status"`
	StatusCode  int        `json:"statusCode" example:"200"`
	Message     string     `json:"message,omitempty"`
	IsArray     bool       `json:"isArray" example:"true"`
	IsPaginated bool       `json:"isPaginated" example:"true"`
	Meta        *MetaModel `json:"paginationMeta,omitempty"`
	Data        any        `json:"data"`
}

type MetaModel struct {
	Total  int `json:"total,omitempty" example:"100"`
	Limit  int `json:"limit,omitempty" example:"10"`
	Offset int `json:"offset,omitempty" example:"0"`
}

type APIResponseError struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// --- Core Response Writers ---

func writeResponseSuccess(w http.ResponseWriter, data APIResponseSuccess) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(data.StatusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func writeResponseError(w http.ResponseWriter, status string, message string, errors []ErrorDetail, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponseError{
		Status:  status,
		Message: message,
		Errors:  errors,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func sendSuccess(w http.ResponseWriter, res APIResponseSuccess) {
	writeResponseSuccess(w, res)
}

func sendError(w http.ResponseWriter, message string, errors []ErrorDetail, statusCode int) {
	writeResponseError(w, "error", message, errors, statusCode)
}

// --- Informational Responses ---

func Send100Continue(w http.ResponseWriter, message string) {
	sendSuccess(w, APIResponseSuccess{Status: "info", StatusCode: http.StatusContinue, Message: message})
}

func Send101SwitchingProtocols(w http.ResponseWriter, message string) {
	sendSuccess(w, APIResponseSuccess{Status: "info", StatusCode: http.StatusSwitchingProtocols, Message: message})
}

func Send102Processing(w http.ResponseWriter, message string) {
	sendSuccess(w, APIResponseSuccess{Status: "info", StatusCode: http.StatusProcessing, Message: message})
}

func Send103EarlyHints(w http.ResponseWriter, message string) {
	sendSuccess(w, APIResponseSuccess{Status: "info", StatusCode: http.StatusEarlyHints, Message: message})
}

// --- Success Responses ---

func Send200OK(w http.ResponseWriter, message string, data any, isArray bool, isPaginated bool, meta *MetaModel) {
	sendSuccess(w, APIResponseSuccess{
		Status:      "success",
		StatusCode:  http.StatusOK,
		Message:     message,
		IsArray:     isArray,
		IsPaginated: isPaginated,
		Meta:        meta,
		Data:        data,
	})
}

func Send201Created(w http.ResponseWriter, message string, data any, isArray bool, isPaginated bool, meta *MetaModel) {
	sendSuccess(w, APIResponseSuccess{
		Status:      "success",
		StatusCode:  http.StatusCreated,
		Message:     message,
		IsArray:     isArray,
		IsPaginated: isPaginated,
		Meta:        meta,
		Data:        data,
	})
}

func Send202Accepted(w http.ResponseWriter, message string, data any, isArray bool, isPaginated bool, meta *MetaModel) {
	sendSuccess(w, APIResponseSuccess{
		Status:      "success",
		StatusCode:  http.StatusAccepted,
		Message:     message,
		IsArray:     isArray,
		IsPaginated: isPaginated,
		Meta:        meta,
		Data:        data,
	})
}

func Send203NonAuthoritativeInfo(w http.ResponseWriter, message string, data any, isArray bool, isPaginated bool, meta *MetaModel) {
	sendSuccess(w, APIResponseSuccess{
		Status:      "success",
		StatusCode:  http.StatusNonAuthoritativeInfo,
		Message:     message,
		IsArray:     isArray,
		IsPaginated: isPaginated,
		Meta:        meta,
		Data:        data,
	})
}

func Send204NoContent(w http.ResponseWriter) {
	sendSuccess(w, APIResponseSuccess{
		Status:     "success",
		StatusCode: http.StatusNoContent,
		Message:    "No Content",
	})
}

func Send205ResetContent(w http.ResponseWriter, message string) {
	sendSuccess(w, APIResponseSuccess{
		Status:     "success",
		StatusCode: http.StatusResetContent,
		Message:    message,
	})
}

func Send206PartialContent(w http.ResponseWriter, message string, data any, isArray bool, isPaginated bool, meta *MetaModel) {
	sendSuccess(w, APIResponseSuccess{
		Status:      "success",
		StatusCode:  http.StatusPartialContent,
		Message:     message,
		IsArray:     isArray,
		IsPaginated: isPaginated,
		Meta:        meta,
		Data:        data,
	})
}

func Send207MultiStatus(w http.ResponseWriter, message string, data any, isArray bool, isPaginated bool, meta *MetaModel) {
	sendSuccess(w, APIResponseSuccess{
		Status:      "success",
		StatusCode:  http.StatusMultiStatus,
		Message:     message,
		IsArray:     isArray,
		IsPaginated: isPaginated,
		Meta:        meta,
		Data:        data,
	})
}

func Send208AlreadyReported(w http.ResponseWriter, message string, data any, isArray bool, isPaginated bool, meta *MetaModel) {
	sendSuccess(w, APIResponseSuccess{
		Status:      "success",
		StatusCode:  http.StatusAlreadyReported,
		Message:     message,
		IsArray:     isArray,
		IsPaginated: isPaginated,
		Meta:        meta,
		Data:        data,
	})
}

// --- Error Responses ---

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

func Send409Conflict(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "conflict", Message: message}}, http.StatusConflict)
}

func Send410Gone(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "resource", Message: message}}, http.StatusGone)
}

func Send415UnsupportedMediaType(w http.ResponseWriter, message string) {
	sendError(w, message, []ErrorDetail{{Field: "media_type", Message: message}}, http.StatusUnsupportedMediaType)
}

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
