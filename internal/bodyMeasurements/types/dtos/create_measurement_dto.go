package bodyMeasurementDtos

import (
	"fmt"
	"reflect"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

/* KG and CM */
type CreateMeasurementDto struct {
	Mass          float64
	Stature       float64
	SittingHeight float64
	ArmSpan       float64

	/* SkinFolds */
	Triceps      float64
	Subscapular  float64
	Biceps       float64
	IliacCrest   float64
	Supraspinale float64
	Abdominal    float64
	FrontThigh   float64
	MedialCalf   float64

	/* Girths */
	Head       float64
	Neck       float64
	ArmRelaxed float64
	ArmFlex    float64
	Forearm    float64
	Wrist      float64
	Chest      float64
	Waist      float64
	Hip        float64
	ThighHigh  float64
	ThighLow   float64
	Calf       float64
	Ankle      float64
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
