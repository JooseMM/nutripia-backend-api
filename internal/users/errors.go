package users

import (
	"fmt"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"net/http"
)

type UserError string

const (
	NOT_FOUND    UserError = "USER_NOT_FOUND"
	EMAIL_IN_USE UserError = "USER_EMAIL_ALREADY_IN_USE"
)

func NotFound(filterTerm string) *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusNotFound,
		ErrorCode:   string(NOT_FOUND),
		Description: fmt.Sprintf("User associated with value '%s' not found", filterTerm),
	}
}

func EmailAlreadyExisting(email string) *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusConflict,
		ErrorCode:   string(EMAIL_IN_USE),
		Description: fmt.Sprintf("Email address '%s' is already in use", email),
	}
}
