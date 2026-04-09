package sessionTypes

import "github.com/google/uuid"

type Session struct {
	ID uuid.UUID
	SessionHash string
	UserId uuid.UUID
}

