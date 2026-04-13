package authenticationDtos

import (
	"strings"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type ChangePasswordRequest struct {
	UserId   uuid.UUID `json:"userId"`
	Password string    `json:"password"`
}

func (d *ChangePasswordRequest) Validate() *core.BaseError {
	var errList []string

	d.Password = strings.TrimSpace(d.Password)
	err := core.ValidatePassword(d.Password)

	errList = append(errList, err...)

	if len(errList) == 0 {
		return nil
	}

	return core.ValidationError(errList)
}
