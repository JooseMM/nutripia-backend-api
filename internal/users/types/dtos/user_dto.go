package userDtos

import (
	"time"
)

type UserDto struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"emailAddress"`
	BirthDate    time.Time `json:"birthDate"`
	Role         Role      `json:"role"`
}
