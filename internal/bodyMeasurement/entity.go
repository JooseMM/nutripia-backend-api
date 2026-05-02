package bodyMeasurement

import (
	bodyMeasurementDtos "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurement/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type BodyMeasurement interface {
	ToDTO() bodyMeasurementDtos.BodyMeasurementDto
	Id() valueobject.Identifier
	ToDB() bodyMeasurementDB
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

func (e *entity) Id() valueobject.Identifier {
	return e.id
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

func (e *entity) ToDTO() bodyMeasurementDtos.BodyMeasurementDto {
	return bodyMeasurementDtos.BodyMeasurementDto{
		Id:            e.id.Key(),
		Mass:          e.mass.Float(),
		Stature:       e.stature.Float(),
		SittingHeight: e.sittingHeight.Float(),
		ArmSpan:       e.armSpan.Float(),

		Triceps:      e.triceps.Float(),
		Subscapular:  e.subscapular.Float(),
		Biceps:       e.biceps.Float(),
		IliacCrest:   e.iliacCrest.Float(),
		Supraspinale: e.supraspinale.Float(),
		Abdominal:    e.abdominal.Float(),
		FrontThigh:   e.frontThigh.Float(),
		MedialCalf:   e.medialCalf.Float(),

		Head:       e.head.Float(),
		Neck:       e.neck.Float(),
		ArmRelaxed: e.armRelaxed.Float(),
		ArmFlex:    e.armFlex.Float(),
		Forearm:    e.forearm.Float(),
		Wrist:      e.wrist.Float(),
		Chest:      e.chest.Float(),
		Waist:      e.waist.Float(),
		Hip:        e.hip.Float(),
		ThighHigh:  e.thighHigh.Float(),
		ThighLow:   e.thighLow.Float(),
		Calf:       e.calf.Float(),
		Ankle:      e.ankle.Float(),

		CreatedAt: e.createdAt.ToTime(),
		ClientId:  e.clientId.Key(),
	}
}

func (e *entity) ToDB() bodyMeasurementDB {
	return bodyMeasurementDB{
		Id:            e.id.Key(),
		Mass:          e.mass.Float(),
		Stature:       e.stature.Float(),
		SittingHeight: e.sittingHeight.Float(),
		ArmSpan:       e.armSpan.Float(),

		Triceps:      e.triceps.Float(),
		Subscapular:  e.subscapular.Float(),
		Biceps:       e.biceps.Float(),
		IliacCrest:   e.iliacCrest.Float(),
		Supraspinale: e.supraspinale.Float(),
		Abdominal:    e.abdominal.Float(),
		FrontThigh:   e.frontThigh.Float(),
		MedialCalf:   e.medialCalf.Float(),

		Head:       e.head.Float(),
		Neck:       e.neck.Float(),
		ArmRelaxed: e.armRelaxed.Float(),
		ArmFlex:    e.armFlex.Float(),
		Forearm:    e.forearm.Float(),
		Wrist:      e.wrist.Float(),
		Chest:      e.chest.Float(),
		Waist:      e.waist.Float(),
		Hip:        e.hip.Float(),
		ThighHigh:  e.thighHigh.Float(),
		ThighLow:   e.thighLow.Float(),
		Calf:       e.calf.Float(),
		Ankle:      e.ankle.Float(),

		CreatedAt: e.createdAt.ToTime(),
		ClientId:  e.clientId.Key(),
	}
}
