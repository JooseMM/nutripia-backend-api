package bodyMeasurementTypes

import (
	"time"

	bodyMeasurementDtos "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements/types/dtos"
	"github.com/google/uuid"
)

/* KG and CM */
type BodyMeasurement struct {
	ID            uuid.UUID
	Mass          float64
	Stature       float64
	SittingHeight float64
	ArmSpan       float64

	/* SkinFolds */
	Triceps      float64
	Subscapular  float64
	Biceps       float64
	IliacCrest   float64
	Supraspinale float64
	Abdominal    float64
	FrontThigh   float64
	MedialCalf   float64

	/* Girths */
	Head       float64
	Neck       float64
	ArmRelaxed float64
	ArmFlex    float64
	Forearm    float64
	Wrist      float64
	Chest      float64
	Waist      float64
	Hip        float64
	ThighHigh  float64
	ThighLow   float64
	Calf       float64
	Ankle      float64

	UserId    uuid.UUID `gorm:"index"`
	CreatedAt time.Time `gorm:"index;type:date"`
}

func (b *BodyMeasurement) ToDTO() *bodyMeasurementDtos.BodyMeasurementDto {
	return &bodyMeasurementDtos.BodyMeasurementDto{
		ID:            b.ID,
		Mass:          b.Mass,
		Stature:       b.Stature,
		SittingHeight: b.SittingHeight,
		ArmSpan:       b.ArmSpan,

		Triceps:       b.Triceps,
		Subscapular:   b.Subscapular,
		Biceps:        b.Biceps,
		IliacCrest:    b.IliacCrest,
		Supraspinale:  b.Supraspinale,
		Abdominal:     b.Abdominal,
		FrontThigh:    b.FrontThigh,
		MedialCalf:    b.MedialCalf,

		Head:          b.Head,
		Neck:          b.Neck,
		ArmRelaxed:    b.ArmRelaxed,
		ArmFlex:       b.ArmFlex,
		Forearm:       b.Forearm,
		Wrist:         b.Wrist,
		Chest:         b.Chest,
		Waist:         b.Waist,
		Hip:           b.Hip,
		ThighHigh:     b.ThighHigh,
		ThighLow:      b.ThighLow,
		Calf:          b.Calf,
		Ankle:         b.Ankle,

		CreatedAt: b.CreatedAt,
	}
}
