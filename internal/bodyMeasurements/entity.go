package bodyMeasurement

import (
	bodyMeasurementDtos "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type BodyMeasurement interface {
	ToDTO() bodyMeasurementDtos.BodyMeasurementDto
	Id() valueobject.Identifier
}

type entity struct {
	id            valueobject.Identifier
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

	clientId  valueobject.Identifier
	createdAt valueobject.Dater
}

func (d *entity) Id() valueobject.Identifier {
	return d.id
}

func NewBodyMeasurement(
	dto bodyMeasurementDtos.CreateMeasurement,
	ownerId valueobject.Identifier,
) BodyMeasurement {
	return &entity{
		id:            valueobject.NewIdentifier(),
		mass:          dto.Mass(),
		stature:       dto.Stature(),
		sittingHeight: dto.SittingHeight(),
		armSpan:       dto.ArmSpan(),

		/* SkinFolds */
		triceps:      dto.Triceps(),
		subscapular:  dto.Subscapular(),
		biceps:       dto.Biceps(),
		iliacCrest:   dto.IliacCrest(),
		supraspinale: dto.Supraspinale(),
		abdominal:    dto.Abdominal(),
		frontThigh:   dto.FrontThigh(),
		medialCalf:   dto.MedialCalf(),

		/* Girths */
		head:       dto.Head(),
		neck:       dto.Neck(),
		armRelaxed: dto.ArmRelaxed(),
		armFlex:    dto.ArmFlex(),
		forearm:    dto.Forearm(),
		wrist:      dto.Wrist(),
		chest:      dto.Chest(),
		waist:      dto.Waist(),
		hip:        dto.Hip(),
		thighHigh:  dto.ThighHigh(),
		thighLow:   dto.ThighLow(),
		calf:       dto.Calf(),
		ankle:      dto.Ankle(),

		/* Metadata */
		clientId:  ownerId,
		createdAt: valueobject.NewTracker(),
	}
}

func (b *entity) ToDTO() bodyMeasurementDtos.BodyMeasurementDto {
	return bodyMeasurementDtos.BodyMeasurementDto{
		Id:            b.id.Key(),
		Mass:          b.mass.Float(),
		Stature:       b.stature.Float(),
		SittingHeight: b.sittingHeight.Float(),
		ArmSpan:       b.armSpan.Float(),

		Triceps:      b.triceps.Float(),
		Subscapular:  b.subscapular.Float(),
		Biceps:       b.biceps.Float(),
		IliacCrest:   b.iliacCrest.Float(),
		Supraspinale: b.supraspinale.Float(),
		Abdominal:    b.abdominal.Float(),
		FrontThigh:   b.frontThigh.Float(),
		MedialCalf:   b.medialCalf.Float(),

		Head:       b.head.Float(),
		Neck:       b.neck.Float(),
		ArmRelaxed: b.armRelaxed.Float(),
		ArmFlex:    b.armFlex.Float(),
		Forearm:    b.forearm.Float(),
		Wrist:      b.wrist.Float(),
		Chest:      b.chest.Float(),
		Waist:      b.waist.Float(),
		Hip:        b.hip.Float(),
		ThighHigh:  b.thighHigh.Float(),
		ThighLow:   b.thighLow.Float(),
		Calf:       b.calf.Float(),
		Ankle:      b.ankle.Float(),

		CreatedAt: b.createdAt.ToTime(),
		ClientId:  b.clientId.Key(),
	}
}
