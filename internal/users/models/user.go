package models

import "github.com/JooseMM/nutripia-backend-api/internal/users/models/value_objects"

type User struct {
	value_objects.UserIdentity
	value_objects.AuthenticationInformation
	value_objects.TrackingInformation
}
