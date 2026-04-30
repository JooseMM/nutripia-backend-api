package valueobject

type MeasureUnit float64

const (
	KG MeasureUnit = iota
	CM
	MM
)

type Measurer interface {
	Float() float64
	Unit() MeasureUnit
	IsEqual(Measurer) bool
}

type measurement struct {
	value float64
	unit  MeasureUnit
}

func (m *measurement) IsEqual(target Measurer) bool {
	if target.Unit() != m.Unit() {
		return false
	}
	return m.Float() == target.Float()
}

func (m *measurement) Float() float64 {
	return m.value
}

func (m *measurement) Unit() MeasureUnit {
	return m.unit
}

func NewMeasurement(raw float64, unitType MeasureUnit) (Measurer, *string) {
	if raw < 0 {
		e := "unit: must not be lower than zero "
		return nil, &e
	}

	return &measurement{
		value: raw,
		unit:  unitType,
	}, nil
}
