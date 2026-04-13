package bodyMeasurementDtos

import (
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type GetByRange struct {
	Start *time.Time
	End   *time.Time
}

func (d *GetByRange) Validate() *core.BaseError {
	var errList []string
	now := time.Now().UTC()

	end := d.End.UTC()
	d.End = &end
	start := d.Start.UTC()
	d.Start = &start

	if d.End.Before(*d.Start) {
		errList = append(
			errList,
			"end: The end date must be later than the start date.")
	}

	if d.End.After(now) {
		errList = append(
			errList,
			"end: The end date cannot be in the future.")
	}

	if len(errList) == 0 {
		return nil
	}
	return core.ValidationError(errList)
}
