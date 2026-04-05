package models

import (
	"strings"
	"time"

	"github.com/JooseMM/nutripia-backend-api/core/errors"
)

type CreateUserDto struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"email_address"`
	Password     string    `json:"password"`
	BirthDate    time.Time `json:"birth_date"`
}

func (d *CreateUserDto) Validate() []*errors.ValidationError {
	var errors []*errors.ValidationError

	d.Firstname = strings.TrimSpace(d.Firstname)
	d.Lastname = strings.TrimSpace(d.Lastname)
	nameErr := d.validateName()
	if nameErr != nil {
		errors = append(errors, nameErr...)
	}

	d.EmailAddress = strings.TrimSpace(d.EmailAddress)
	emailErr := d.validateEmailAddress()
	if emailErr != nil {
		errors = append(errors, emailErr...)
	}

	d.Password = strings.TrimSpace(d.Password)
	passErr := d.validatePassword()
	if passErr != nil {
		errors = append(errors, passErr)
	}

	return errors
}

func (d *CreateUserDto) validateEmailAddress() []*errors.ValidationError {

	return nil
}

func (d *CreateUserDto) validateName() []*errors.ValidationError {
	var errorList []*errors.ValidationError

	if d.Firstname == "" {
		errorList = append(errorList, &errors.ValidationError{
			Field:   "firstname",
			Message: "is required",
		})
	}

	if d.Lastname == "" {
		errorList = append(errorList, &errors.ValidationError{
			Field:   "lastname",
			Message: "is required",
		})
	}

	return errorList
}

func (d *CreateUserDto) validatePassword() *errors.ValidationError {
	if d.Password == "" {
		return &errors.ValidationError{
			Field:   "lastname",
			Message: "is required",
		}
	}
	return nil
}
