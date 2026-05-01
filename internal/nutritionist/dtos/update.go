package nutritionistDtos

import (
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type RawUpdateNutritionistRequest struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"emailAddress"`
	BirthDate    time.Time `json:"birthDate"`
}

type UpdateNutritionist struct {
	name         valueobject.Namer
	emailAddress valueobject.Emailer
	birthDate    valueobject.Dater
}

func (u *UpdateNutritionist) BirthDate() valueobject.Dater {
	return u.birthDate
}

func (u *UpdateNutritionist) Name() valueobject.Namer {
	return u.name
}

func (u *UpdateNutritionist) EmailAddress() valueobject.Emailer {
	return u.emailAddress
}

func (dto *RawUpdateNutritionistRequest) ToValueObject() (*UpdateNutritionist, *core.BaseError) {
	var errList []string

	name, nameErr := valueobject.NewName(dto.Firstname, dto.Lastname)
	if nameErr != nil {
		errList = append(errList, nameErr.Details...)
	}

	email, emailErr := valueobject.NewEmailAddress(dto.EmailAddress)
	if emailErr != nil {
		errList = append(errList, emailErr.Details...)
	}

	birthDate, birthDateErr := valueobject.NewBirthDate(dto.BirthDate)
	if birthDateErr != nil {
		errList = append(errList, birthDateErr.Details...)
	}

	if len(errList) > 0 {
		return nil, core.ValidationError(errList)
	}

	return &UpdateNutritionist{
		name:         name,
		emailAddress: email,
		birthDate:    birthDate,
	}, nil
}
