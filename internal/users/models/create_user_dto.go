package users

import (
	"strings"
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type CreateUserRequest struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"emailAddress"`
	Password     string    `json:"password"`
	BirthDate    time.Time `json:"birthDate"`
}

type CreateUserResponse struct {
	UserId string `json:"userId"`
}

func (d *CreateUserRequest) Validate() *core.BaseError {
	var errorList []string

	d.Firstname = strings.TrimSpace(d.Firstname)
	d.Lastname = strings.TrimSpace(d.Lastname)

	errorList = append(errorList, d.validateName()...)

	d.EmailAddress = strings.TrimSpace(d.EmailAddress)
	emailErr := d.validateEmailAddress()
	if emailErr != nil {
		errorList = append(errorList, *emailErr)
	}

	d.Password = strings.TrimSpace(d.Password)
	passErr := d.validatePassword()
	if passErr != nil {
		errorList = append(errorList, *passErr)
	}

	if len(errorList) == 0 {
		return nil
	}

	return core.ValidationError(errorList)
}

func (d *CreateUserRequest) validateEmailAddress() *string {
	return nil
}

func (d *CreateUserRequest) validateName() []string {
	var errList []string

	if d.Firstname == "" {
		errList = append(errList, "firstname: is required")
	}

	if d.Lastname == "" {
		errList = append(errList, "lastname: is required")
	}

	return errList
}

func (d *CreateUserRequest) validatePassword() *string {
	if d.Password == "" {
		var description = "lastname: is required"
		return &description
	}

	return nil
}
