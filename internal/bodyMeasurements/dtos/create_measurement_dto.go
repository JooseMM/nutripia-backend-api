package bodyMeasurementDtos

import (
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type RawCreateMeasurement struct {
	Mass          float64 `json:"mass"`
	Stature       float64 `json:"stature"`
	SittingHeight float64 `json:"sittingHeight"`
	ArmSpan       float64 `json:"armSpan"`

	/* SkinFolds */
	Triceps      float64 `json:"triceps"`
	Subscapular  float64 `json:"subscapular"`
	Biceps       float64 `json:"biceps"`
	IliacCrest   float64 `json:"iliacCrest"`
	Supraspinale float64 `json:"supraspinale"`
	Abdominal    float64 `json:"abdominal"`
	FrontThigh   float64 `json:"frontThigh"`
	MedialCalf   float64 `json:"medialCalf"`

	/* Girths */
	Head       float64 `json:"head"`
	Neck       float64 `json:"neck"`
	ArmRelaxed float64 `json:"armRelaxed"`
	ArmFlex    float64 `json:"armFlex"`
	Forearm    float64 `json:"forearm"`
	Wrist      float64 `json:"wrist"`
	Chest      float64 `json:"chest"`
	Waist      float64 `json:"waist"`
	Hip        float64 `json:"hip"`
	ThighHigh  float64 `json:"thighHigh"`
	ThighLow   float64 `json:"thighLow"`
	Calf       float64 `json:"calf"`
	Ankle      float64 `json:"ankle"`
}

func (d *RawCreateMeasurement) ToValueObject() (*CreateMeasurement, *core.BaseError) {
	var errList []string
	// Helper to reduce boilerplate
	validate := func(val float64, unit valueobject.MeasureUnit) valueobject.Measurer {
		m, err := valueobject.NewMeasurement(val, unit)
		if err != nil {
			errList = append(errList, err.Details...)
		}
		return m
	}

	// Basic Stats
	mass := validate(d.Mass, valueobject.KG)
	stature := validate(d.Stature, valueobject.CM)
	sittingHeight := validate(d.SittingHeight, valueobject.CM)
	armSpan := validate(d.ArmSpan, valueobject.CM)

	/* SkinFolds - Standardized to MM usually, change to CM if preferred */
	triceps := validate(d.Triceps, valueobject.MM)
	subscapular := validate(d.Subscapular, valueobject.MM)
	biceps := validate(d.Biceps, valueobject.MM)
	iliacCrest := validate(d.IliacCrest, valueobject.MM)
	supraspinale := validate(d.Supraspinale, valueobject.MM)
	abdominal := validate(d.Abdominal, valueobject.MM)
	frontThigh := validate(d.FrontThigh, valueobject.MM)
	medialCalf := validate(d.MedialCalf, valueobject.MM)

	/* Girths */
	head := validate(d.Head, valueobject.CM)
	neck := validate(d.Neck, valueobject.CM)
	armRelaxed := validate(d.ArmRelaxed, valueobject.CM)
	armFlex := validate(d.ArmFlex, valueobject.CM)
	forearm := validate(d.Forearm, valueobject.CM)
	wrist := validate(d.Wrist, valueobject.CM)
	chest := validate(d.Chest, valueobject.CM)
	waist := validate(d.Waist, valueobject.CM)
	hip := validate(d.Hip, valueobject.CM)
	thighHigh := validate(d.ThighHigh, valueobject.CM)
	thighLow := validate(d.ThighLow, valueobject.CM)
	calf := validate(d.Calf, valueobject.CM)
	ankle := validate(d.Ankle, valueobject.CM)

	if len(errList) > 0 {
		return nil, core.ValidationError(errList)
	}

	return &CreateMeasurement{
		mass:          mass,
		stature:       stature,
		sittingHeight: sittingHeight,
		armSpan:       armSpan,
		triceps:       triceps,
		subscapular:   subscapular,
		biceps:        biceps,
		iliacCrest:    iliacCrest,
		supraspinale:  supraspinale,
		abdominal:     abdominal,
		frontThigh:    frontThigh,
		medialCalf:    medialCalf,
		head:          head,
		neck:          neck,
		armRelaxed:    armRelaxed,
		armFlex:       armFlex,
		forearm:       forearm,
		wrist:         wrist,
		chest:         chest,
		waist:         waist,
		hip:           hip,
		thighHigh:     thighHigh,
		thighLow:      thighLow,
		calf:          calf,
		ankle:         ankle,
	}, nil
}

type CreateMeasurement struct {
	mass          valueobject.Measurer
	stature       valueobject.Measurer
	sittingHeight valueobject.Measurer
	armSpan       valueobject.Measurer

	/* SkinFolds */
	triceps      valueobject.Measurer
	subscapular  valueobject.Measurer
	biceps       valueobject.Measurer
	iliacCrest   valueobject.Measurer
	supraspinale valueobject.Measurer
	abdominal    valueobject.Measurer
	frontThigh   valueobject.Measurer
	medialCalf   valueobject.Measurer

	/* Girths */
	head       valueobject.Measurer
	neck       valueobject.Measurer
	armRelaxed valueobject.Measurer
	armFlex    valueobject.Measurer
	forearm    valueobject.Measurer
	wrist      valueobject.Measurer
	chest      valueobject.Measurer
	waist      valueobject.Measurer
	hip        valueobject.Measurer
	thighHigh  valueobject.Measurer
	thighLow   valueobject.Measurer
	calf       valueobject.Measurer
	ankle      valueobject.Measurer
}

func (d *CreateMeasurement) Mass() valueobject.Measurer {
	return d.mass
}

func (m *CreateMeasurement) Stature() valueobject.Measurer {
	return m.stature
}

func (m *CreateMeasurement) SittingHeight() valueobject.Measurer {
	return m.sittingHeight
}

func (m *CreateMeasurement) ArmSpan() valueobject.Measurer {
	return m.armSpan
}

// --- Skinfolds ---

func (m *CreateMeasurement) Triceps() valueobject.Measurer {
	return m.triceps
}

func (m *CreateMeasurement) Subscapular() valueobject.Measurer {
	return m.subscapular
}

func (m *CreateMeasurement) Biceps() valueobject.Measurer {
	return m.biceps
}

func (m *CreateMeasurement) IliacCrest() valueobject.Measurer {
	return m.iliacCrest
}

func (m *CreateMeasurement) Supraspinale() valueobject.Measurer {
	return m.supraspinale
}

func (m *CreateMeasurement) Abdominal() valueobject.Measurer {
	return m.abdominal
}

func (m *CreateMeasurement) FrontThigh() valueobject.Measurer {
	return m.frontThigh
}

func (m *CreateMeasurement) MedialCalf() valueobject.Measurer {
	return m.medialCalf
}

// --- Girths ---

func (m *CreateMeasurement) Head() valueobject.Measurer {
	return m.head
}

func (m *CreateMeasurement) Neck() valueobject.Measurer {
	return m.neck
}

func (m *CreateMeasurement) ArmRelaxed() valueobject.Measurer {
	return m.armRelaxed
}

func (m *CreateMeasurement) ArmFlex() valueobject.Measurer {
	return m.armFlex
}

func (m *CreateMeasurement) Forearm() valueobject.Measurer {
	return m.forearm
}

func (m *CreateMeasurement) Wrist() valueobject.Measurer {
	return m.wrist
}

func (m *CreateMeasurement) Chest() valueobject.Measurer {
	return m.chest
}

func (m *CreateMeasurement) Waist() valueobject.Measurer {
	return m.waist
}

func (m *CreateMeasurement) Hip() valueobject.Measurer {
	return m.hip
}

func (m *CreateMeasurement) ThighHigh() valueobject.Measurer {
	return m.thighHigh
}

func (m *CreateMeasurement) ThighLow() valueobject.Measurer {
	return m.thighLow
}

func (m *CreateMeasurement) Calf() valueobject.Measurer {
	return m.calf
}

func (m *CreateMeasurement) Ankle() valueobject.Measurer {
	return m.ankle
}
