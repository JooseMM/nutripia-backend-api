package authenticationDtos

import (
	authenticationTypes "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/types"
	"github.com/google/uuid"
)

type LoginResponse struct {
	UserId    uuid.UUID                     `json:"userId"`
	Firstname string                        `json:"firstname"`
	Role      authenticationTypes.UserRoles `json:"role"`
	Token     string                        `json:"token"`
}
