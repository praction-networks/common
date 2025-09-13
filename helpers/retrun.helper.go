package helpers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/praction-networks/common/appError"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/response"
)

func HandleAppError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*appError.AppError); ok {
		// Extract field information using the improved extraction function
		field := extractFieldFromError(appErr)

		errResponse := response.ErrorDetail{
			Field:   field,
			Message: appErr.Message,
		}
		response.SendCustomResponse(w, []response.ErrorDetail{errResponse}, appErr.HTTPCode)
	} else {
		response.Send500InternalServerError(w, "Internal Server Error")
	}
}

// extractFieldFromError safely extracts field information from error messages
func extractFieldFromError(err *appError.AppError) string {
	// Try to extract field from message pattern like "field: message"
	if strings.Contains(err.Message, ": ") {
		parts := strings.SplitN(err.Message, ": ", 2)
		if len(parts) == 2 && parts[0] != "" {
			return parts[0]
		}
	}

	// Fallback to error code for non-validation errors
	return string(err.Code)
}

func ParseRequestBodyAndRespond(r *http.Request, w http.ResponseWriter, dest any) bool {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		logger.Error("Error parsing request body", err)
		response.Send400BadRequest(w, "Invalid request body: "+err.Error())
		return false
	}
	return true
}

func ValidateRequestAndRespond(w http.ResponseWriter, validationErrors []response.ErrorDetail, logMessage string) bool {
	if len(validationErrors) > 0 {
		logger.Warn(logMessage, "validationErrors", validationErrors)
		response.SendCustomError(w, "Validation failed", validationErrors, http.StatusBadRequest)
		return false
	}
	return true
}
