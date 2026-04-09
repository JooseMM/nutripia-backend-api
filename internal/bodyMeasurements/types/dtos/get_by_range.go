package bodyMeasurementDtos

import "time"

type GetByRange struct {
	Start *time.Time
	End   *time.Time
}
