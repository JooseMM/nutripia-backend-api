package verificationTypes

import (
	"time"

	"github.com/google/uuid"
)

type VerificationToken struct {
	ID        uuid.UUID
	Token     string
	UserId    uuid.UUID
	CreatedAt time.Time
	ExpiredAt time.Time
}
