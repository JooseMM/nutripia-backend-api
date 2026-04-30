package clients

import (
	"time"

	clientDtos "github.com/JooseMM/nutripia-backend-api/internal/clients/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type entity struct {
	id        valueobject.Identifier
	name      valueobject.Namer
	email     valueobject.Emailer
	birthDate valueobject.BirthDater
	ownerId   valueobject.Identifier

	createdAt time.Time
	updateAt  time.Time
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
		birthDate: e.birthDate.Date(),
		createdAt: e.createdAt,
		updatedAt: e.updateAt,
	}
}

func (e *entity) Id() valueobject.Identifier {
	return e.id
}

func FromDTO(dto clientDtos.CreateClientRequest, ownerId valueobject.Identifier) Client {
	now := time.Now().UTC()
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

	e.updateAt = time.Now().UTC()
}

func (e *entity) ToDTO() clientDtos.ClientDto {
	return clientDtos.ClientDto{
		ID:           e.id.Key(),
		Firstname:    e.name.Firstname(),
		Lastname:     e.name.Lastname(),
		EmailAddress: e.email.String(),
		BirthDate:    e.birthDate.Date(),
	}
}
