package valueobject

import "github.com/JooseMM/nutripia-backend-api/pkg/core"

type RUTer interface {
	ToString() string
}

type rut struct {
	value string
}

type RutRequest struct {
	Value string
}

func (dto *RutRequest) ToValueObject() (RUTer, *core.BaseError) {
	if dto.Value == "" {
		return nil, core.ValidationError([]string{"RUT: invalid rut format"})
	}

	return &rut{
		value: dto.Value,
	}, nil
}

func (r *rut) ToString() string {
	return r.value
}

func NewRUT(raw string) (RUTer, *core.BaseError) {
	if raw == "" {
		return nil, core.ValidationError([]string{"RUT: is required"})
	}
	return &rut{value: raw}, nil
}
