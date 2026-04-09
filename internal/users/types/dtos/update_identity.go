package userDtos

import (
	"strings"
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type UpdateIdentityRequest struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"emailAddress"`
	BirthDate    time.Time `json:"birthDate"`
}

func (d *UpdateIdentityRequest) Validate() *core.BaseError {
	var errorList []string

	d.Firstname = strings.TrimSpace(d.Firstname)
	d.Lastname = strings.TrimSpace(d.Lastname)

	errorList = append(errorList, ValidateName(d.Firstname, d.Lastname)...)

	d.EmailAddress = strings.TrimSpace(d.EmailAddress)
	emailErr := ValidateEmailAddress(d.EmailAddress)
	if emailErr != nil {
		errorList = append(errorList, *emailErr)
	}

	if len(errorList) == 0 {
		return nil
	}

	return core.ValidationError(errorList)
}
