package nutritionistDtos

import (
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type RegisterNutritionistInput struct {
	Fullname     valueobject.FullnameRequest `json:"fullname"`
	EmailAddress valueobject.EmailRequest    `json:"emailAddress"`
	Password     valueobject.PasswordRequest `json:"password"`
	BirthDate    valueobject.DateRequest     `json:"birthDate"`
	Rut          valueobject.RutRequest      `json:"rut"`
}

func (dto *RegisterNutritionistInput) ToEntity() (*RegisterNutritionist, *core.BaseError) {
	var errList []string

	fullname, err := dto.Fullname.ToValueObject()
	if err != nil {
		errList = append(errList, err.Details...)
	}

	email, err := dto.EmailAddress.ToValueObject()
	if err != nil {
		errList = append(errList, err.Details...)
	}

	password, err := dto.Password.ToValueObject()
	if err != nil {
		errList = append(errList, err.Details...)
	}

	rut, err := dto.Rut.ToValueObject()
	if err != nil {
		errList = append(errList, err.Details...)
	}

	return &RegisterNutritionist{
		fullName:     fullname,
		emailAddress: email,
		password:     password,
		birthDate:    dto.BirthDate.ToValueObject(),
		rut:          rut,
	}, nil
}

type RegisterNutritionist struct {
	fullName     valueobject.Namer
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
	return dto.fullName
}

func (dto *RegisterNutritionist) BirthDate() valueobject.Dater {
	return dto.birthDate
}

func (dto *RegisterNutritionist) RUT() valueobject.RUTer {
	return dto.rut
}
