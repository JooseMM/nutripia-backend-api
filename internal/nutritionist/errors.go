package nutritionist

import (
	"fmt"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"net/http"
)

type NutritionistError string

const (
	NUTRITIONIST_NOT_FOUND           NutritionistError = "NUTRITIONIST_NOT_FOUND"
	NUTRITIONIST_EMAIL_ALREADY_IN_US NutritionistError = "NUTRITIONIST_EMAIL_ALREADY_IN_USE"
)

func UserNotFound(filterTerm string) *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusNotFound,
		ErrorCode:   string(NUTRITIONIST_NOT_FOUND),
		Description: fmt.Sprintf("Nutritionist associated with value '%s' not found", filterTerm),
	}
}

func UserEmailAlreadyExisting(email string) *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusConflict,
		ErrorCode:   string(NUTRITIONIST_EMAIL_ALREADY_IN_US),
		Description: fmt.Sprintf("Email address '%s' is already in use", email),
	}
}
