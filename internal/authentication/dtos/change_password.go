package authenticationDtos

import (
	"github.com/google/uuid"
)

type ChangePasswordRequest struct {
	UserId   uuid.UUID `json:"userId"`
	Password string    `json:"password"`
}
