package valueobject

import (
	"strings"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type Emailer interface {
	String() string
}

type EmailAddress struct {
	prefix string
	domain string
}

func (e *EmailAddress) String() string {
	return e.prefix + "@" + e.domain
}

func NewEmailAddress(raw string) (Emailer, *core.BaseError) {
	var errList []string
	email := strings.TrimSpace(raw)

	if email == "" {
		errList = append(errList, "EmailAddress: is required")
	}

	parts := strings.Split(email, "@")
	if len(parts) < 2 {
		errList = append(errList, "EmailAddress: wrong email address format")
	}

	if len(errList) > 0 {
		return nil, core.ValidationError(errList)
	}

	return &EmailAddress{
		prefix: parts[0],
		domain: parts[1],
	}, nil
}
