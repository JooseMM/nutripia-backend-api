package nutritionist

import (
	nutritionistDtos "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type Nutritionist interface {
	Id() valueobject.Identifier
	Name() valueobject.Namer
	EmailAddress() valueobject.Emailer
	BirthDate() valueobject.Dater

	Password() valueobject.Passworder
	UpdatePassword(password valueobject.Passworder)

	RUT() valueobject.RUTer
	IsEmailConfirmed() bool
	MarkAsEmailAddressConfirmed()

	Update(data nutritionistDtos.UpdateNutritionist)
	ToDTO() nutritionistDtos.Nutritionist
}

type Entity struct {
	id               valueobject.Identifier
	name             valueobject.Namer
	isEmailConfirmed bool
	emailAddress     valueobject.Emailer
	birthDate        valueobject.Dater
	password         valueobject.Passworder
	rut              valueobject.RUTer
	createdAt        valueobject.Dater
	updatedAt        valueobject.Dater
}

func (e *Entity) Id() valueobject.Identifier {
	return e.id
}

func (e *Entity) Name() valueobject.Namer {
	return e.name
}

func (e *Entity) EmailAddress() valueobject.Emailer {
	return e.emailAddress
}

func (e *Entity) Password() valueobject.Passworder {
	return e.password
}

func (e *Entity) UpdatePassword(password valueobject.Passworder) {
	e.password = password
	e.updatedAt = valueobject.NewTracker()
}

func (e *Entity) RUT() valueobject.RUTer {
	return e.rut
}

func (e *Entity) BirthDate() valueobject.Dater {
	return e.birthDate
}

func (e *Entity) IsEmailConfirmed() bool {
	return e.isEmailConfirmed
}

func (e *Entity) MarkAsEmailAddressConfirmed() {
	e.isEmailConfirmed = true
	e.updatedAt = valueobject.NewTracker()
}

func (e *Entity) Update(update nutritionistDtos.UpdateNutritionist) {
	e.name = update.Name()
	e.emailAddress = update.EmailAddress()
	e.birthDate = update.BirthDate()
	e.updatedAt = valueobject.NewTracker()
}

func (e *Entity) ToDTO() nutritionistDtos.Nutritionist {
	return nutritionistDtos.Nutritionist{
		ID:               e.id.Key(),
		Firstname:        e.name.Firstname(),
		Lastname:         e.name.Lastname(),
		EmailAddress:     e.EmailAddress().String(),
		IsEmailConfirmed: e.isEmailConfirmed,
		BirthDate:        e.birthDate.ToTime(),
	}
}

func FromRegisterDTO(data nutritionistDtos.RegisterNutritionist) Nutritionist {
	now := valueobject.NewTracker()
	return &Entity{
		id:               valueobject.NewIdentifier(),
		name:             data.Name(),
		emailAddress:     data.Email(),
		isEmailConfirmed: false,
		birthDate:        data.BirthDate(),
		password:         data.Password(),
		rut:              data.RUT(),
		createdAt:        now,
		updatedAt:        now,
	}
}

func NewEntity(
	name valueobject.Namer,
	emailAddress valueobject.Emailer,
	isEmailConfirmed bool,
	birthDate valueobject.Dater,
	password valueobject.Passworder,
	rut valueobject.RUTer,
) Nutritionist {
	now := valueobject.NewTracker()
	return &Entity{
		id:               valueobject.NewIdentifier(),
		name:             name,
		emailAddress:     emailAddress,
		isEmailConfirmed: isEmailConfirmed,
		birthDate:        birthDate,
		password:         password,
		rut:              rut,
		createdAt:        now,
		updatedAt:        now,
	}
}
