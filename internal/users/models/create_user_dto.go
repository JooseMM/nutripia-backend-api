package models

import (
	"strings"
	"time"
)

type CreateUserDto struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"email_address"`
	Password     string    `json:"password"`
	BirthDate    time.Time `json:"birth_date"`
}

func (d *CreateUserDto) Validate() *string {
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

	finalMsg := strings.Join(errorList, ",\n")
	return &finalMsg
}

func (d *CreateUserDto) validateEmailAddress() *string {
	return nil
}

func (d *CreateUserDto) validateName() []string {
	var errList []string

	if d.Firstname == "" {
		errList = append(errList, "firstname: is required")
	}

	if d.Lastname == "" {
		errList = append(errList, "lastname: is required")
	}

	return errList
}

func (d *CreateUserDto) validatePassword() *string {
	if d.Password == "" {
		var description = "lastname: is required"
		return &description
	}

	return nil
}
