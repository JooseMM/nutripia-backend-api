package clientDtos

import (
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type RawUpdateClientRequest struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"emailAddress"`
	BirthDate    time.Time `json:"birthDate"`
}

func (dto *RawUpdateClientRequest) ToValueObject() (*UpdateClientRequest, *core.BaseError) {
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

	return &UpdateClientRequest{
		Name:         name,
		EmailAddress: email,
		BirthDate:    birthDate,
	}, nil
}

type UpdateClientRequest struct {
	Name         valueobject.Namer
	EmailAddress valueobject.Emailer
	BirthDate    valueobject.Dater
}
