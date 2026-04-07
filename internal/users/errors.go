package users

import (
	"fmt"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

func NotFound(id uuid.UUID) *core.BaseError {
	return &core.BaseError{
		ErrorCode:   "USER_NOT_FOUND",
		Description: "User '" + id.String() + "' not found",
	}
}

func EmailAlreadyExisting(email string) *core.BaseError {
	return &core.BaseError{
		ErrorCode:   "USER_EMAIL_ALREADY_IN_USE",
		Description: fmt.Sprintf("Email address %s is already in use", email),
	}
}
