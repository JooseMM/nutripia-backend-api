package session

import (
	"time"

	authenticationTypes "github.com/JooseMM/nutripia-backend-api/internal/security/authentication/types"
	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID
	SessionHash string
	UserId      uuid.UUID
	Role        authenticationTypes.UserRoles
	ExpiredAt   time.Time
}
