package clientTypes

import (
	"github.com/google/uuid"
)

type Client struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	ClientIdentity
	TrackingInformation
}
