package clientDtos

import (
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type RawCreateClientRequest struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"emailAddress"`
	BirthDate    time.Time `json:"birthDate"`
}

type CreateClientRequest struct {
	Name         valueobject.Namer
	EmailAddress valueobject.Emailer
	BirthDate    valueobject.Dater
}

func (dto *RawCreateClientRequest) ToValueObject() (*CreateClientRequest, *core.BaseError) {
	var errList []string
	name, err := valueobject.NewName(dto.Firstname, dto.Lastname)
	if err != nil {
		errList = append(errList, err...)
	}

	email, err := valueobject.NewEmailAddress(dto.EmailAddress)
	if err != nil {
		errList = append(errList, err...)
	}

	birthDate, err := valueobject.NewBirthDate(dto.BirthDate)
	if err != nil {
		errList = append(errList, err...)
	}

	if len(errList) > 0 {
		return nil, core.ValidationError(errList)
	}

	return &CreateClientRequest{
		Name:         name,
		EmailAddress: email,
		BirthDate:    birthDate,
	}, nil
}
