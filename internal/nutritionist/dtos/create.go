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

func (dto *RawRegisterNutritionist) ToValueObject() (*RegisterNutritionist, *core.BaseError) {
	var errList []string

	name, nameErr := valueobject.NewName(dto.Firstname, dto.Lastname)
	if nameErr != nil {
		errList = append(errList, nameErr.Details...)
	}

	email, emailErr := valueobject.NewEmailAddress(dto.EmailAddress)
	if emailErr != nil {
		errList = append(errList, emailErr.Details...)
	}

	password, passwordErr := valueobject.NewPassword(dto.Password)
	if passwordErr != nil {
		errList = append(errList, passwordErr.Details...)
	}

	birthDate, birthDateErr := valueobject.NewBirthDate(dto.BirthDate)
	if birthDateErr != nil {
		errList = append(errList, birthDateErr.Details...)
	}

	rut, rutErr := valueobject.NewRUT(dto.RUT)
	if rutErr != nil {
		errList = append(errList, rutErr.Details...)
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

type RegisterNutritionist struct {
	name         valueobject.Namer
	emailAddress valueobject.Emailer
	password     valueobject.Passworder
	birthDate    valueobject.Dater
	rut          valueobject.RUTer
}

func (dto *RegisterNutritionist) Email() valueobject.Emailer {
	return dto.emailAddress
}

func (dto *RegisterNutritionist) Password() valueobject.Passworder {
	return dto.password
}

func (dto *RegisterNutritionist) Name() valueobject.Namer {
	return dto.name
}

func (dto *RegisterNutritionist) BirthDate() valueobject.Dater {
	return dto.birthDate
}

func (dto *RegisterNutritionist) RUT() valueobject.RUTer {
	return dto.rut
}
