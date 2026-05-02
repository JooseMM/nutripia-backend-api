package valueobject

import (
	"time"

	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type DateType int

const (
	Appointment DateType = iota
	BirthDate
	Tracker
	RangeDate
	Expiration
)

type Dater interface {
	ToTime() time.Time
	YearsSince() int
	Type() DateType
	IsEqual(target Dater) bool
	IsPastOrNow() bool
}

type date struct {
	value    time.Time
	dateType DateType
}

type DateRequest struct {
	Value    time.Time
	DateType DateType
}

func (dto *DateRequest) ToValueObject() Dater {
	return &date{
		value:    dto.Value,
		dateType: dto.DateType,
	}
}

func (d *date) Type() DateType {
	return d.dateType
}

func (d *date) ToTime() time.Time {
	return d.value
}

func (d *date) IsPastOrNow() bool {
	return isPastOrNow(d.value)
}

func (d *date) IsEqual(target Dater) bool {
	if d.dateType != target.Type() {
		return false
	}

	return d.value.Equal(target.ToTime())
}

func (d *date) YearsSince() int {
	now := time.Now().UTC()
	years := now.Year() - d.value.Year()

	if now.Before(d.value.AddDate(years, 0, 0)) {
		years--
	}

	return years
}

func NewExpiredDate(minutesFromNow int) Dater {
	return &date{
		value:    time.Now().UTC().Add(time.Duration(minutesFromNow) * time.Minute),
		dateType: Expiration,
	}
}

func ExpiredDateFromTime(raw time.Time) (Dater, *core.BaseError) {
	return &date{
		value:    raw.UTC(),
		dateType: Expiration,
	}, nil
}

func NewAppointmentDate(raw time.Time) (Dater, *core.BaseError) {
	var errList []string

	if isPastOrNow(raw) {
		errList = append(errList, "AppointmentDate: cannot be today or in the past.")
	}

	return &date{
		value:    raw.UTC(),
		dateType: Appointment,
	}, nil
}

func NewBirthDate(raw time.Time) (Dater, *core.BaseError) {
	if !isPastOrNow(raw) {
		return nil, core.ValidationError([]string{"BirthDate: cannot be today or in the future."})
	}

	return &date{
		value:    raw.UTC(),
		dateType: BirthDate,
	}, nil
}

func NewTracker() Dater {
	return &date{
		value:    time.Now().UTC(),
		dateType: Tracker,
	}
}

func NewRangeDate(raw time.Time) (Dater, *core.BaseError) {
	if !isPastOrNow(raw) {
		return nil, core.ValidationError([]string{"RageDate: cannot be today or in the future."})
	}

	return &date{
		value:    raw.UTC(),
		dateType: RangeDate,
	}, nil
}

func NewTrackerFromTime(raw *time.Time) (Dater, *core.BaseError) {
	now := time.Now().UTC()

	if raw.UTC().After(now) {
		return nil, core.ValidationError([]string{"Tracker: cannot be in the future."})
	}
	return &date{
		value:    raw.UTC(),
		dateType: Tracker,
	}, nil
}

func isPastOrNow(raw time.Time) bool {
	now := time.Now().UTC()
	rawUTC := raw.UTC()

	return rawUTC.Before(now) || rawUTC.Equal(now)
}
