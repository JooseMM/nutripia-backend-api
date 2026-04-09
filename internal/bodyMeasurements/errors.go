package bodyMeasurement

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type BodyMeasurementError string

const (
	BODY_MEASUREMENT_RECORD_NOT_FOUND       BodyMeasurementError = "BODY_MEASUREMENT_RECORD_NOT_FOUND"
	BODY_MEASUREMENT_RECORD_ALREADY_CREATED BodyMeasurementError = "BODY_MEASUREMENT_RECORD_ALREADY_CREATED"
)

func BodyMeasurementNotFound(filterTerm string) *core.BaseError {
	return &core.BaseError{
		StatusCode:  http.StatusNotFound,
		ErrorCode:   string(BODY_MEASUREMENT_RECORD_NOT_FOUND),
		Description: fmt.Sprintf("Body measurement associated with value '%s' not found", filterTerm),
	}
}

func BodyMeasurementAlreadyExisting(date time.Time) *core.BaseError {
	return &core.BaseError{
		StatusCode: http.StatusConflict,
		ErrorCode:  string(BODY_MEASUREMENT_RECORD_ALREADY_CREATED),
		Description: fmt.Sprintf(
			"Body measurement for date '%s' is already present",
			date.Format("2000-01-01"),
		),
	}
}
