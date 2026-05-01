package valueobject

import (
	"strings"
	"unicode/utf8"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type Namer interface {
	ToString() string
	Firstname() string
	Lastname() string
}

type HumanName struct {
	firstname string
	lastname  string
}

func (n *HumanName) Firstname() string {
	return n.firstname
}

func (n *HumanName) Lastname() string {
	return n.lastname
}

func (n *HumanName) ToString() string {
	return n.firstname + " " + n.lastname
}

func NewName(rawFirstname string, rawLastname string) (Namer, *core.BaseError) {
	var errList []string
	firstname := strings.TrimSpace(rawFirstname)
	lastname := strings.TrimSpace(rawLastname)

	fCount := utf8.RuneCountInString(firstname)
	lCount := utf8.RuneCountInString(lastname)

	if firstname == "" {
		errList = append(errList, "Firstname: is required")
	}
	if fCount < 3 {
		errList = append(errList, "Firstname: must be equal or greater than 3 characters")
	}
	if core.HasHarmfulSymbols(firstname) {
		errList = append(
			errList,
			"Firstname: Invalid characters detected. Please use only letters and standard punctuation.",
		)
	}

	if lastname == "" {
		errList = append(errList, "Lastname: is required")
	}
	if lCount < 3 {
		errList = append(errList, "Lastname: must be equal or greater than 3 characters")
	}
	if core.HasHarmfulSymbols(lastname) {
		errList = append(
			errList,
			"Lastname: Invalid characters detected. Please use only letters and standard punctuation.",
		)
	}

	if len(errList) > 0 {
		return nil, core.ValidationError(errList)
	}

	return &HumanName{
		firstname,
		lastname,
	}, nil
}
