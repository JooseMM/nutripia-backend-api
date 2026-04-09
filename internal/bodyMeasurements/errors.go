package bodyMeasurement

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type BodyMeasurementError string

const (
	BODY_MEASUREMENT_RECORD_NOT_FOUND       BodyMeasurementError = "BODY_MEASUREMENT_RECORD_NOT_FOUND"
	BODY_MEASUREMENT_RECORD_ALREADY_CREATED BodyMeasurementError = "BODY_MEASUREMENT_RECORD_ALREADY_CREATED"
	USER_MUST_BE_CLIENT                     BodyMeasurementError = "USER_MUST_BE_CLIENT"
)

func BodyMeasurementNotFound(filterTerm string) *core.BaseError {
	return &core.BaseError{
		StatusCode: http.StatusNotFound,
		ErrorCode:  string(BODY_MEASUREMENT_RECORD_NOT_FOUND),
		Description: fmt.Sprintf(
			"Body measurement associated with value '%s' not found",
			filterTerm,
		),
	}
}

func BodyMeasurementNotAClient(userId uuid.UUID) *core.BaseError {
	return &core.BaseError{
		StatusCode: http.StatusUnprocessableEntity,
		ErrorCode:  string(USER_MUST_BE_CLIENT),
		Description: fmt.Sprintf(
			"Body measurements can only be created for users with a 'client' role. (User ID: %s)",
			userId.String(),
		),
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
