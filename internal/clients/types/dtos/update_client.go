package clientDtos

import (
	"strings"
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type UpdateClientRequest struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"emailAddress"`
	BirthDate    time.Time `json:"birthDate"`
}

func (d *UpdateClientRequest) Validate() *core.BaseError {
	var errorList []string

	d.Firstname = strings.TrimSpace(d.Firstname)
	d.Lastname = strings.TrimSpace(d.Lastname)

	errorList = append(errorList, core.ValidateName(d.Firstname, d.Lastname)...)

	d.EmailAddress = strings.TrimSpace(d.EmailAddress)
	emailErr := core.ValidateEmailAddress(d.EmailAddress)
	errorList = append(errorList, emailErr...)

	if len(errorList) == 0 {
		return nil
	}

	return core.ValidationError(errorList)
}
