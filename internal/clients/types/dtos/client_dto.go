package clientDtos

import (
	"github.com/google/uuid"
	"time"
)

type ClientDto struct {
	ID           uuid.UUID `json:"id"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"emailAddress"`
	BirthDate    time.Time `json:"birthDate"`
}
