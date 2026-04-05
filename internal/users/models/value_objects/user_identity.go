package value_objects

import (
	"time"

	"github.com/google/uuid"
)

type UserIdentity struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Firstname    string    `gorm:"size:100;not null"`
	Lastname     string    `gorm:"size:100;not null"`
	EmailAddress string    `gorm:"uniqueIndex;size:255;not null"`
	DateBirth    time.Time `gorm:"type:date"`
}
