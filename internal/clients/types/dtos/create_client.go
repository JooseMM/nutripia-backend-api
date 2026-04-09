package clientDtos

import (
	"strings"
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type CreateClientRequest struct {
	Firstname           string    `json:"firstname"`
	Lastname            string    `json:"lastname"`
	EmailAddress        string    `json:"emailAddress"`
	BirthDate           time.Time `json:"birthDate"`
	NutritionistOwnerId uuid.UUID `json:"nutritionistOwnerId"`
}

func (d *CreateClientRequest) Validate() *core.BaseError {
	var errorList []string

	d.Firstname = strings.TrimSpace(d.Firstname)
	d.Lastname = strings.TrimSpace(d.Lastname)

	errorList = append(errorList, ValidateName(d.Firstname, d.Lastname)...)

	d.EmailAddress = strings.TrimSpace(d.EmailAddress)
	emailErr := ValidateEmailAddress(d.EmailAddress)
	if emailErr != nil {
		errorList = append(errorList, *emailErr)
	}

	if len(errorList) == 0 {
		return nil
	}

	return core.ValidationError(errorList)
}
