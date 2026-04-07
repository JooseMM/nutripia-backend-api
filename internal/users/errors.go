package users

import (
	"fmt"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"net/http"
)

func NotFound(filterTerm string) *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusNotFound,
		ErrorCode:   "USER_NOT_FOUND",
		Description: fmt.Sprintf("User associated with value '%s' not found", filterTerm),
	}
}

func EmailAlreadyExisting(email string) *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusConflict,
		ErrorCode:   "USER_EMAIL_ALREADY_IN_USE",
		Description: fmt.Sprintf("Email address '%s' is already in use", email),
	}
}
