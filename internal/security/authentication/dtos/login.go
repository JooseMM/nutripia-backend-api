package authenticationDtos

import (
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/google/uuid"
)

type RawLoginRequest struct {
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
}

func (r *RawLoginRequest) ToValueObject() (*LoginRequest, *core.BaseError) {
	var errList []string
	email, err := valueobject.NewEmailAddress(r.EmailAddress)
	errList = append(errList, err...)

	if r.Password == "" {
		errList = append(errList, "password: is required")
	}

	if len(errList) > 0 {
		return nil, core.ValidationError(errList)
	}

	return &LoginRequest{
		rawPassword:  r.Password,
		emailAddress: email,
	}, nil
}

type LoginRequest struct {
	emailAddress valueobject.Emailer
	rawPassword  string
}

func (l *LoginRequest) EmailAddress() valueobject.Emailer {
	return l.emailAddress
}

func (l *LoginRequest) Password() string {
	return l.rawPassword
}

type LoginResponse struct {
	UserId    uuid.UUID             `json:"userId"`
	Firstname string                `json:"firstname"`
	Role      valueobject.UserRoles `json:"role"`
	Token     string                `json:"token"`
}

func NewLoginResponse(
	userId valueobject.Identifier,
	userName valueobject.Namer,
	role valueobject.UserRoles,
	token valueobject.Tokenizer,
) LoginResponse {
	return LoginResponse{
		UserId:    userId.Key(),
		Firstname: userName.Firstname(),
		Role:      role,
		Token:     token.String(),
	}
}
