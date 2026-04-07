package userModels

import (
	"strings"
	"time"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type UserIdentityDto struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"emailAddress"`
	Password     string    `json:"password"`
	BirthDate    time.Time `json:"birthDate"`
}

func (d *UserIdentityDto) Validate() *core.BaseError {
	var errorList []string

	d.Firstname = strings.TrimSpace(d.Firstname)
	d.Lastname = strings.TrimSpace(d.Lastname)

	errorList = append(errorList, ValidateName(d.Firstname, d.Lastname)...)

	d.EmailAddress = strings.TrimSpace(d.EmailAddress)
	emailErr := ValidateEmailAddress(d.EmailAddress)
	if emailErr != nil {
		errorList = append(errorList, *emailErr)
	}

	d.Password = strings.TrimSpace(d.Password)
	passErr := ValidatePassword(d.Password)
	if passErr != nil {
		errorList = append(errorList, *passErr)
	}

	if len(errorList) == 0 {
		return nil
	}

	return core.ValidationError(errorList)
}
