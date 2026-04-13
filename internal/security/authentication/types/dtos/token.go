package authenticationDtos

import "github.com/JooseMM/nutripia-backend-api/pkg/core"

type Token struct {
	Token string `json:"token"`
}

func (d *Token) Validate() *core.BaseError {
	var errList []string

	if d.Token == "" {
		errList = append(errList, "token: is required")
	}

	if len(errList) == 0 {
		return nil
	}

	return core.ValidationError(errList)
}
