package iamguard

import "errors"

var (
	errMissingService = errors.New("iamguard: Config.Service is required")
	errMissingRouter  = errors.New("iamguard: Config.Router is required")
)
