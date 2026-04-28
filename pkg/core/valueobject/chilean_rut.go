package valueobject

type RUTer interface {
	ToString() string
}

type Rut struct {
	value string
}

func (r *Rut) ToString() string {
	return r.value
}

func NewRUT(raw string) (RUTer, []string) {
	return &Rut{value: raw}, nil
}
