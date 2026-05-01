package valueobject

import "github.com/JooseMM/nutripia-backend-api/pkg/core"

type RUTer interface {
	ToString() string
}

type Rut struct {
	value string
}

func (r *Rut) ToString() string {
	return r.value
}

func NewRUT(raw string) (RUTer, *core.BaseError) {
	return &Rut{value: raw}, nil
}
