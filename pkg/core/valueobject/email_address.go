package valueobject

import (
	"strings"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type Emailer interface {
	String() string
}

type EmailRequest struct {
	Prefix string `json:"prefix"`
	Domain string `json:"domain"`
}

func (dto *EmailRequest) ToValueObject() (Emailer, *core.BaseError) {
	if dto.Prefix == "" {
		return nil, core.ValidationError([]string{"Prefix: is required"})
	}

	if dto.Domain == "" {
		return nil, core.ValidationError([]string{"Domain: is required"})
	}

	return &emailAddress{
		prefix: dto.Prefix,
		domain: dto.Domain,
	}, nil
}

type emailAddress struct {
	prefix string
	domain string
}

func (e *emailAddress) String() string {
	return e.prefix + "@" + e.domain
}

func NewEmailAddress(raw string) (Emailer, *core.BaseError) {
	if raw == "" {
		return nil, core.ValidationError([]string{"EmailAddress: is required"})
	}

	parts := strings.Split(raw, "@")
	return &emailAddress{
		prefix: parts[0],
		domain: parts[1],
	}, nil
}
