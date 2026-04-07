package users

import "github.com/JooseMM/nutripia-backend-api/internal/users/models/value_objects"

type User struct {
	users.UserIdentity
	users.AuthenticationInformation
	users.TrackingInformation
}
