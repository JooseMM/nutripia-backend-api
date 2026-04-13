package nutritionistDtos

import (
	"strings"
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type CreateNutritionistRequest struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	EmailAddress string    `json:"emailAddress"`
	Password     string    `json:"password"`
	BirthDate    time.Time `json:"birthDate"`
	RUT          string    `json:"rut"`
}

func (d *CreateNutritionistRequest) Validate() *core.BaseError {
	var errorList []string

	d.Firstname = strings.TrimSpace(d.Firstname)
	d.Lastname = strings.TrimSpace(d.Lastname)
	errorList = append(errorList, core.ValidateName(d.Firstname, d.Lastname)...)

	d.EmailAddress = strings.TrimSpace(d.EmailAddress)
	errorList = append(errorList, core.ValidateEmailAddress(d.EmailAddress)...)

	d.Password = strings.TrimSpace(d.Password)
	errorList = append(errorList, core.ValidatePassword(d.Password)...)

	d.BirthDate = d.BirthDate.UTC()
	errorList = append(errorList, core.ValidateBirthDate(d.BirthDate)...)

	d.RUT = strings.TrimSpace(d.RUT)
	errorList = append(errorList, core.ValidateRUT(&d.RUT)...)

	if len(errorList) == 0 {
		return nil
	}

	return core.ValidationError(errorList)
}
