package clients

import (
	clientDtos "github.com/JooseMM/nutripia-backend-api/internal/clients/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type entity struct {
	id        valueobject.Identifier
	name      valueobject.Namer
	email     valueobject.Emailer
	birthDate valueobject.Dater
	ownerId   valueobject.Identifier

	createdAt valueobject.Dater
	updateAt  valueobject.Dater
}

type Client interface {
	Id() valueobject.Identifier
	Update(dto clientDtos.UpdateClientRequest)
	ToDB() entityDB
	ToDTO() clientDtos.ClientDto
}

func (e *entity) ToDB() entityDB {
	return entityDB{
		id:        e.id.Key(),
		ownerId:   e.ownerId.Key(),
		firstname: e.name.Firstname(),
		lastname:  e.name.Lastname(),
		email:     e.email.String(),
		birthDate: e.birthDate.ToTime(),
		createdAt: e.createdAt.ToTime(),
		updatedAt: e.updateAt.ToTime(),
	}
}

func (e *entity) Id() valueobject.Identifier {
	return e.id
}

func NewClient(dto clientDtos.CreateClientRequest, ownerId valueobject.Identifier) Client {
	now := valueobject.NewTracker()
	return &entity{id: valueobject.NewIdentifier(),
		name:      dto.Name,
		email:     dto.EmailAddress,
		birthDate: dto.BirthDate,
		ownerId:   ownerId,
		createdAt: now,
		updateAt:  now,
	}
}

func (e *entity) Update(dto clientDtos.UpdateClientRequest) {
	e.name = dto.Name
	e.email = dto.EmailAddress
	e.birthDate = dto.BirthDate

	e.updateAt = valueobject.NewTracker()
}

func (e *entity) ToDTO() clientDtos.ClientDto {
	return clientDtos.ClientDto{
		ID:           e.id.Key(),
		Firstname:    e.name.Firstname(),
		Lastname:     e.name.Lastname(),
		EmailAddress: e.email.String(),
		BirthDate:    e.birthDate.ToTime(),
	}
}
