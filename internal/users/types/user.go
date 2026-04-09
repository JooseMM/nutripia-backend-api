package userTypes

import (
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	UserIdentity
	AuthenticationInformation

	TrackingInformation
}
