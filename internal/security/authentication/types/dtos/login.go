package authenticationDtos

import (
	"strings"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type LoginRequest struct {
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
}

func (d *LoginRequest) Validate() *core.BaseError {
	var errList []string

	d.Password = strings.TrimSpace(d.Password)
	errList = append(errList, core.ValidatePassword(d.Password)...)

	d.EmailAddress = strings.TrimSpace(d.EmailAddress)
	errList = append(errList, core.ValidateEmailAddress(d.EmailAddress)...)

	if len(errList) == 0 {
		return nil
	}

	return core.ValidationError(errList)
}
