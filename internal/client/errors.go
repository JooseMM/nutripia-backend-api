package clients

import (
	"fmt"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"net/http"
)

type ClientError string

const (
	CLIENT_NOT_FOUND    ClientError = "CLIENT_NOT_FOUND"
	CLIENT_EMAIL_IN_USE ClientError = "CLIENT_EMAIL_ALREADY_IN_USE"
)

func UserNotFound(filterTerm string) *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusNotFound,
		ErrorCode:   string(CLIENT_NOT_FOUND),
		Description: fmt.Sprintf("Client associated with value '%s' not found", filterTerm),
	}
}

func UserEmailAlreadyExisting(email string) *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusConflict,
		ErrorCode:   string(CLIENT_EMAIL_IN_USE),
		Description: fmt.Sprintf("Email address '%s' is already in use", email),
	}
}
