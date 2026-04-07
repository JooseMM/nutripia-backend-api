package core

import "net/http"

type BaseError struct {
	StatusCode  uint16   `json:"statusCode"`
	ErrorCode   string   `json:"errorCode"`
	Description string   `json:"description"`
	Details     []string `json:"details,omitempty"`
}

func UnexpectedError(description string) *BaseError {
	return &BaseError{
		StatusCode:  http.StatusInternalServerError,
		ErrorCode:   "UNEXPECTED_ERROR",
		Description: description,
		Details:     nil,
	}
}

func ValidationError(detailList []string) *BaseError {
	return &BaseError{
		StatusCode:  http.StatusUnprocessableEntity,
		ErrorCode:   "VALIDATION_ERROR",
		Description: "The request body contains invalid data or failed business logic validation.",
		Details:     detailList,
	}
}
