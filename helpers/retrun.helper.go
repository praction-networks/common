package helpers

import (
	"encoding/json"
	"net/http"

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
