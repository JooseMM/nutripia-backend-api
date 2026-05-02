package authentication

import (
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type AuthenticationError string

const (
	WRONG_CREDENTIALS   AuthenticationError = "WRONG_CREDENTIALS"
	EMAIL_NOT_CONFIRMED AuthenticationError = "EMAIL_NOT_CONFIRMED"
)

func WrongCredentials() *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusUnauthorized,
		ErrorCode:   string(WRONG_CREDENTIALS),
		Description: "Invalid credentials",
	}
}

func EmailNotConfirmed() *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusForbidden,
		ErrorCode:   string(EMAIL_NOT_CONFIRMED),
		Description: "User haven't confirmed it's email address yet.",
	}
}
