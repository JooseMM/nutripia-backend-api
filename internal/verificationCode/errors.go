package verification

import (
	"net/http"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type VerificationCodeError string

const (
	VERIFICATION_CODE_NOT_FOUND = "VERIFICATION_CODE_NOT_FOUND"
	VERIFICATION_CODE_EXPIRED   = "VERIFICATION_CODE_EXPIRED"
)

func VerificationCodeNotFound() *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusNotFound,
		ErrorCode:   string(VERIFICATION_CODE_NOT_FOUND),
		Description: "Session not found",
	}
}

func VerificationCodeExpired() *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusForbidden,
		ErrorCode:   string(VERIFICATION_CODE_EXPIRED),
		Description: "Session expired",
	}
}
