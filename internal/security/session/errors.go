package session

import (
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type SessionError string

const (
	SESSION_NOT_FOUND = "SESSION_NOT_FOUND"
	SESSION_EXPIRED   = "SESSION_EXPIRED"
)

func SessionNotFound() *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusNotFound,
		ErrorCode:   string(SESSION_NOT_FOUND),
		Description: "Session not found",
	}
}

func SessionExpired() *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusForbidden,
		ErrorCode:   string(SESSION_NOT_FOUND),
		Description: "Session expired",
	}
}
