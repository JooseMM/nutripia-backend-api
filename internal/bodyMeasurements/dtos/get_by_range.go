package bodyMeasurementDtos

import (
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

const layout = "2000-01-01"

type RangeDate struct {
	start valueobject.Dater
	end   valueobject.Dater
}

func (r *RangeDate) Start() valueobject.Dater {
	return r.start
}

func (r *RangeDate) End() valueobject.Dater {
	return r.end
}

func NewRangeDate(rawStart string, rawEnd string) (*RangeDate, *core.BaseError) {
	var errList []string
	start, err := time.Parse(layout, rawStart)
	if err != nil {
		start = time.Now().AddDate(0, -1, 0)
	}

	end, err := time.Parse(layout, rawEnd)
	if err != nil {
		end = time.Now()
	}

	if end.Before(start) || end.Equal(start) {
		errList = append(
			errList,
			"Invalid range: the start date must be earlier than the end date.",
		)
	}

	startDater, e := valueobject.NewRangeDate(start)
	errList = append(errList, e...)

	endDater, e := valueobject.NewRangeDate(end)
	errList = append(errList, e...)

	if len(errList) > 0 {
		return nil, core.ValidationError(errList)
	}

	return &RangeDate{
		start: startDater,
		end:   endDater,
	}, nil
}
