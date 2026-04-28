package nutritionistDtos

import (
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type RawRegisterNutritionist struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"emailAddress"`
	Password     string    `json:"password"`
	BirthDate    time.Time `json:"birthDate"`
	RUT          string    `json:"rut"`
}

type RegisterNutritionist struct {
	name         valueobject.Namer
	emailAddress valueobject.Emailer
	password     valueobject.Passworder
	birthDate    valueobject.BirthDater
	rut          valueobject.RUTer
}

func (dto *RawRegisterNutritionist) ToValueObject() (*RegisterNutritionist, *core.BaseError) {
	var errList []string

	name, nameErr := valueobject.NewName(dto.Firstname, dto.Lastname)
	if nameErr != nil {
		errList = append(errList, nameErr...)
	}

	email, emailErr := valueobject.NewEmailAddress(dto.EmailAddress)
	if emailErr != nil {
		errList = append(errList, emailErr...)
	}

	password, passwordErr := valueobject.NewPassword(dto.Password)
	if passwordErr != nil {
		errList = append(errList, passwordErr...)
	}

	birthDate, birthDateErr := valueobject.NewBirthDate(dto.BirthDate)
	if birthDateErr != nil {
		errList = append(errList, birthDateErr...)
	}

	rut, rutErr := valueobject.NewRUT(dto.RUT)
	if rutErr != nil {
		errList = append(errList, rutErr...)
	}

	if len(errList) > 0 {
		return nil, core.ValidationError(errList)
	}

	return &RegisterNutritionist{
		name:         name,
		emailAddress: email,
		password:     password,
		birthDate:    birthDate,
		rut:          rut,
	}, nil
}
