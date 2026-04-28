package nutritionist

import (
	"time"

	nutritionistDtos "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type Nutritionist interface {
	Id() valueobject.Identifier
	Name() valueobject.Namer
	EmailAddress() valueobject.Emailer
	BirthDate() valueobject.BirthDater
	Password() valueobject.Passworder
	RUT() valueobject.RUTer

	Update(data nutritionistDtos.UpdateNutritionist)
	ToDTO() nutritionistDtos.Nutritionist
}

type Entity struct {
	id           valueobject.Identifier
	name         valueobject.Namer
	emailAddress valueobject.Emailer
	birthDate    valueobject.BirthDater
	password     valueobject.Passworder
	rut          valueobject.RUTer

	createdAt time.Time
	updatedAt time.Time
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

func (e *Entity) RUT() valueobject.RUTer {
	return e.rut
}

func (e *Entity) BirthDate() valueobject.BirthDater {
	return e.birthDate
}

func (e *Entity) Update(update nutritionistDtos.UpdateNutritionist) {
	e.name = update.Name()
	e.emailAddress = update.EmailAddress()
	e.birthDate = update.BirthDate()
	e.updatedAt = time.Now()
}

func (e *Entity) ToDTO() nutritionistDtos.Nutritionist {
	return nutritionistDtos.Nutritionist{
		ID:           e.id.Key(),
		Firstname:    e.name.Firstname(),
		Lastname:     e.name.Lastname(),
		EmailAddress: e.EmailAddress().String(),
		BirthDate:    e.birthDate.Date(),
	}
}

func NewEntity(
	name valueobject.Namer,
	emailAddress valueobject.Emailer,
	birthDate valueobject.BirthDater,
	password valueobject.Passworder,
	rut valueobject.RUTer,
) Nutritionist {
	now := time.Now()
	return &Entity{
		id:           valueobject.NewIdentifier(),
		name:         name,
		emailAddress: emailAddress,
		birthDate:    birthDate,
		password:     password,
		rut:          rut,
		createdAt:    now,
		updatedAt:    now,
	}
}
