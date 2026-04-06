package response

import "github.com/JooseMM/nutripia-backend-api/pkg/errors"

type ErrorResponse struct {
	Title       string            `json:"title"`
	StatusCode  uint16            `json:"statusCode"`
	ErrorList   []errors.BaseError `json:"error"`
	Description string            `json:"description"`
}
