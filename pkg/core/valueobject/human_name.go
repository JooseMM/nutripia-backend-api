package valueobject

import "github.com/JooseMM/nutripia-backend-api/pkg/core"

type Namer interface {
	ToString() string
	Firstname() string
	Lastname() string
}

type FullnameRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (dto *FullnameRequest) ToValueObject() (Namer, *core.BaseError) {
	if dto.Firstname == "" {
		return nil, core.ValidationError([]string{"Firstname: is required"})
	}

	if dto.Lastname == "" {
		return nil, core.ValidationError([]string{"Lastname: is required"})
	}

	return &fullname{
		firstname: dto.Firstname,
		lastname:  dto.Lastname,
	}, nil
}

type fullname struct {
	firstname string
	lastname  string
}

func (n *fullname) Firstname() string {
	return n.firstname
}

func (n *fullname) Lastname() string {
	return n.lastname
}

func (n *fullname) ToString() string {
	return n.firstname + " " + n.lastname
}

func NewFullName(rawFirstname string, rawLastname string) (Namer, *core.BaseError) {
	var errList []string

	if rawFirstname == "" {
		errList = append(errList, "Firstname: is required")
	}

	if rawLastname == "" {
		errList = append(errList, "Lastname: is required")
	}

	if len(errList) > 0 {
		return nil, core.ValidationError(errList)
	}

	return &fullname{
		rawFirstname,
		rawLastname,
	}, nil
}
