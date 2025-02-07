package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/lucsky/cuid"
	"github.com/praction-networks/common/appError"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/response"
)

func HandleAppError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*appError.AppError); ok {
		errResponse := response.ErrorDetail{
			Field:   string(appErr.Code),
			Message: appErr.Message,
		}
		response.SendCustomResponse(w, []response.ErrorDetail{errResponse}, appErr.HTTPCode)
	} else {
		response.Send500InternalServerError(w, "Internal Server Error")
	}
}

func ParseRequestBodyAndRespond(r *http.Request, w http.ResponseWriter, dest interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		logger.Error("Error parsing request body", err)
		response.Send400BadRequest(w, "Invalid request body: "+err.Error())
		return false
	}
	return true
}

func IsValidCUID(id string, w http.ResponseWriter) bool {
	if err := cuid.IsCuid(id); err != nil {
		logger.Error("Request ID should be a valid CUID", err)
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
