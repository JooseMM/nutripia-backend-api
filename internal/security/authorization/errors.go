package authorization

import (
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type AuthorizationError string

const (
	UNAUTHORIZED AuthorizationError = "UNAUTHORIZED"
	FORBIDDEN    AuthorizationError = "FORBIDDEN"
)

func NotEnoughPermissions() *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusForbidden,
		ErrorCode:   string(FORBIDDEN),
		Description: "User does not have the neccessary permissions",
	}
}

func Unauthenticated() *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusUnauthorized,
		ErrorCode:   string(UNAUTHORIZED),
		Description: "User is not authenticated",
	}
}
