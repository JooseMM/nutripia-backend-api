package bodyMeasurement

import (
	"context"
	"errors"
	"time"

	bodyMeasurementDtos "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type entityDB struct {
	id            uuid.UUID `gorm:"type:uuid;primaryKey"`
	mass          float64   `gorm:"type:decimal(10,2)"`
	stature       float64   `gorm:"type:decimal(10,2)"`
	sittingHeight float64   `gorm:"column:sitting_height;type:decimal(10,2)"`
	armSpan       float64   `gorm:"column:arm_span;type:decimal(10,2)"`

	// Skinfolds
	triceps      float64 `gorm:"type:decimal(10,2)"`
	subscapular  float64 `gorm:"type:decimal(10,2)"`
	biceps       float64 `gorm:"type:decimal(10,2)"`
	iliacCrest   float64 `gorm:"type:decimal(10,2)"`
	supraspinale float64 `gorm:"type:decimal(10,2)"`
	abdominal    float64 `gorm:"type:decimal(10,2)"`
	frontThigh   float64 `gorm:"column:front_thigh;type:decimal(10,2)"`
	medialCalf   float64 `gorm:"column:Medial_calf;type:decimal(10,2)"`

	// Girths/Circumferences
	head       float64 `gorm:"type:decimal(10,2)"`
	neck       float64 `gorm:"type:decimal(10,2)"`
	armRelaxed float64 `gorm:"column:arm_relaxed;type:decimal(10,2)"`
	armFlex    float64 `gorm:"column:arm_flex;type:decimal(10,2)"`
	forearm    float64 `gorm:"type:decimal(10,2)"`
	wrist      float64 `gorm:"type:decimal(10,2)"`
	chest      float64 `gorm:"type:decimal(10,2)"`
	waist      float64 `gorm:"type:decimal(10,2)"`
	hip        float64 `gorm:"type:decimal(10,2)"`
	thighHigh  float64 `gorm:"column:thigh_high;type:decimal(10,2)"`
	thighLow   float64 `gorm:"column:thigh_low;type:decimal(10,2)"`
	calf       float64 `gorm:"type:decimal(10,2)"`
	ankle      float64 `gorm:"type:decimal(10,2)"`

	// Metadata & Relations
	createdAt time.Time `gorm:"column:created_at"`
	updatedAt time.Time `gorm:"column:updated_at"`

	clientId uuid.UUID `gorm:"type:uuid;not null;index"`
}

func (e *entityDB) toEntity() (BodyMeasurement, *core.BaseError) {
	var errList []string

	validate := func(val float64, unit valueobject.MeasureUnit) valueobject.Measurer {
		m, err := valueobject.NewMeasurement(val, unit)
		if err != nil {
			errList = append(errList, *err)
		}
		return m
	}

	// Basic Stats
	mass := validate(e.mass, valueobject.KG)
	stature := validate(e.stature, valueobject.CM)
	sittingHeight := validate(e.sittingHeight, valueobject.CM)
	armSpan := validate(e.armSpan, valueobject.CM)

	/* SkinFolds - Standardized to MM usually, change to CM if preferred */
	triceps := validate(e.triceps, valueobject.MM)
	subscapular := validate(e.subscapular, valueobject.MM)
	biceps := validate(e.biceps, valueobject.MM)
	iliacCrest := validate(e.iliacCrest, valueobject.MM)
	supraspinale := validate(e.supraspinale, valueobject.MM)
	abdominal := validate(e.abdominal, valueobject.MM)
	frontThigh := validate(e.frontThigh, valueobject.MM)
	medialCalf := validate(e.medialCalf, valueobject.MM)

	/* Girths */
	head := validate(e.head, valueobject.CM)
	neck := validate(e.neck, valueobject.CM)
	armRelaxed := validate(e.armRelaxed, valueobject.CM)
	armFlex := validate(e.armFlex, valueobject.CM)
	forearm := validate(e.forearm, valueobject.CM)
	wrist := validate(e.wrist, valueobject.CM)
	chest := validate(e.chest, valueobject.CM)
	waist := validate(e.waist, valueobject.CM)
	hip := validate(e.hip, valueobject.CM)
	thighHigh := validate(e.thighHigh, valueobject.CM)
	thighLow := validate(e.thighLow, valueobject.CM)
	calf := validate(e.calf, valueobject.CM)
	ankle := validate(e.ankle, valueobject.CM)

	if len(errList) > 0 {
		return nil, core.CorrupetedDatabase(errList)
	}

	return &entity{
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

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetAllByUser(
		ctx context.Context,
		userId valueobject.Identifier,
	) ([]BodyMeasurement, *core.BaseError)
	GetByDate(
		ctx context.Context,
		date valueobject.Dater,
		userId valueobject.Identifier,
	) (BodyMeasurement, *core.BaseError)
	DateAlreadyFill(
		ctx context.Context,
		date valueobject.Dater,
		userId valueobject.Identifier,
	) (*bool, *core.BaseError)
	GetByRange(
		ctx context.Context,
		rageDate bodyMeasurementDtos.RangeDate,
		clientId valueobject.Identifier,
	) ([]BodyMeasurement, *core.BaseError)
	GetById(
		ctx context.Context,
		id valueobject.Identifier,
	) (BodyMeasurement, *core.BaseError)
	Create(ctx context.Context, u BodyMeasurement) *core.BaseError
	Delete(ctx context.Context, id valueobject.Identifier) *core.BaseError
}

func NewBodyMeasurementRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetAllByUser(
	ctx context.Context,
	userID valueobject.Identifier,
) ([]BodyMeasurement, *core.BaseError) {
	var rawList []entityDB

	result := r.db.WithContext(ctx).Where("user_id = ?", userID.Key()).Find(&rawList)
	if result.Error != nil {
		return nil, core.UnexpectedError(result.Error.Error())
	}

	var list []BodyMeasurement
	for _, raw := range rawList {
		client, err := raw.toEntity()
		if err != nil {
			return nil, err
		}

		list = append(list, client)
	}

	return list, nil
}

func (r *repository) DateAlreadyFill(
	ctx context.Context,
	date valueobject.Dater,
	userId valueobject.Identifier,
) (*bool, *core.BaseError) {
	var count int64

	result := r.db.WithContext(ctx).
		Model(&entityDB{}).
		Where("created_at::date = ?::date AND client_id = ?", date.ToTime(), userId.Key()).
		Count(&count)

	if result.Error != nil {
		return nil, core.UnexpectedError(result.Error.Error())
	}

	isFound := count > 0
	return &isFound, nil
}

func (r *repository) GetByDate(
	ctx context.Context,
	date valueobject.Dater,
	userId valueobject.Identifier,
) (BodyMeasurement, *core.BaseError) {
	var raw entityDB

	result := r.db.WithContext(ctx).
		First(&raw, "created_at = ? AND user_id = ?", date.ToTime(), userId.Key())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, BodyMeasurementNotFound(date.ToTime().Format("2000-01-01"))
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	measurement, err := raw.toEntity()
	if err != nil {
		return nil, err
	}

	return measurement, nil
}

func (r *repository) GetByRange(
	ctx context.Context,
	rangeDate bodyMeasurementDtos.RangeDate,
	clientId valueobject.Identifier,
) ([]BodyMeasurement, *core.BaseError) {
	var rawList []entityDB

	result := r.db.WithContext(ctx).
		Where("user_id = ?", clientId.Key()).
		Where("created_at BETWEEN ? AND ?", rangeDate.Start().ToTime(), rangeDate.End().ToTime()).
		Order("created_at ASC").
		Find(&rawList)

	if result.Error != nil {
		return nil, core.UnexpectedError(result.Error.Error())
	}

	var measurementList []BodyMeasurement
	for _, raw := range rawList {
		m, err := raw.toEntity()
		if err != nil {
			return nil, err
		}

		measurementList = append(measurementList, m)
	}

	return measurementList, nil
}

func (r *repository) GetById(
	ctx context.Context,
	id valueobject.Identifier,
) (BodyMeasurement, *core.BaseError) {
	var raw entityDB

	result := r.db.WithContext(ctx).First(&raw, "id = ?", id.Key())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, BodyMeasurementNotFound(id.String())
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	measurement, err := raw.toEntity()
	if err != nil {
		return nil, err
	}

	return measurement, nil
}

func (r *repository) Create(
	ctx context.Context,
	measurement BodyMeasurement,
) *core.BaseError {
	result := r.db.WithContext(ctx).Create(measurement.ToDTO())
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id valueobject.Identifier) *core.BaseError {
	result := r.db.WithContext(ctx).Delete(&entityDB{}, id.Key())
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return BodyMeasurementNotFound(id.String())
	}

	return nil
}
