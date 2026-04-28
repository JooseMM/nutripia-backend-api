package valueobject

import (
	"time"
)

type BirthDater interface {
	Date() time.Time
	isUnderAge() bool
	Age() int
}

type BirthDate struct {
	date time.Time
}

func (b *BirthDate) Date() time.Time {
	return b.date
}

func (b *BirthDate) isUnderAge() bool {
	years18Later := b.date.AddDate(18, 0, 0)
	now := time.Now()

	return now.After(years18Later) || now.Equal(years18Later)
}

func (b *BirthDate) Age() int {
	now := time.Now()
	years := now.Year() - b.date.Year()

	if now.Before(b.date.AddDate(years, 0, 0)) {
		years--
	}

	return years
}

func NewBirthDate(raw time.Time) (BirthDater, []string) {
	var errList []string
	now := time.Now()

	if now.Before(raw) || now.Equal(raw) {
		errList = append(errList, "BirthDate: cannot be today or in the future.")
	}

	if len(errList) > 0 {
		return nil, errList
	}

	return &BirthDate{
		date: raw,
	}, nil
}
