package errors

import (
	"fmt"

	"github.com/JooseMM/nutripia-backend-api/pkg/errors"
	"github.com/google/uuid"
)

func NotFound(id uuid.UUID) *errors.BaseError {
	return &errors.BaseError{
		ErrorCode:   "USER_NOT_FOUND",
		Description: "User '" + id.String() + "' not found",
	}
}

func EmailAreadyExisting(email string) *errors.BaseError {
	return &errors.BaseError{
		ErrorCode:   "USER_EMAIL_ALREADY_IN_USE",
		Description: fmt.Sprintf("Email address %s is already in use", email),
	}
}
