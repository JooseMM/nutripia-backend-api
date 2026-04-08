package userTypes

import (
	"time"
)

type UserIdentity struct {
	Firstname    string    `gorm:"size:100;not null"`
	Lastname     string    `gorm:"size:100;not null"`
	EmailAddress string    `gorm:"uniqueIndex;size:255;not null"`
	DateBirth    time.Time `gorm:"type:date"`
}
