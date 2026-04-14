package bodyMeasurementDtos

import (
	"fmt"
	"reflect"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

/* KG and CM */
type CreateMeasurementDto struct {
	ID            uuid.UUID `json:"id"`
	Mass          float64   `json:"mass"`
	Stature       float64   `json:"stature"`
	SittingHeight float64   `json:"sittingHeight"`
	ArmSpan       float64   `json:"armSpan"`

	/* SkinFolds */
	Triceps      float64 `json:"triceps"`
	Subscapular  float64 `json:"subscapular"`
	Biceps       float64 `json:"biceps"`
	IliacCrest   float64 `json:"iliacCrest"`
	Supraspinale float64 `json:"supraspinale"`
	Abdominal    float64 `json:"abdominal"`
	FrontThigh   float64 `json:"frontThigh"`
	MedialCalf   float64 `json:"medialCalf"`

	/* Girths */
	Head       float64 `json:"head"`
	Neck       float64 `json:"neck"`
	ArmRelaxed float64 `json:"armRelaxed"`
	ArmFlex    float64 `json:"armFlex"`
	Forearm    float64 `json:"forearm"`
	Wrist      float64 `json:"wrist"`
	Chest      float64 `json:"chest"`
	Waist      float64 `json:"waist"`
	Hip        float64 `json:"hip"`
	ThighHigh  float64 `json:"thighHigh"`
	ThighLow   float64 `json:"thighLow"`
	Calf       float64 `json:"calf"`
	Ankle      float64 `json:"ankle"`

	ClientId  uuid.UUID `json:"clientId"`
}

func (b *CreateMeasurementDto) Validate() *core.BaseError {
	var errorList []string
	v := reflect.ValueOf(b)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name

		// We only care about float64 fields for these measurements
		if field.Kind() == reflect.Float64 {
			val := field.Float()

			// Shared rule: ISAK measurements cannot be zero or negative
			if val <= 0 {
				var errDetail = fmt.Sprintf("measurement '%s' must be a positive value (got %v)", fieldName, val)
				errorList = append(errorList, errDetail)
			}
		}
	}

	if len(errorList) == 0 {
		return nil
	}

	return core.ValidationError(errorList)
}
