package bodyMeasurement

import (
	"context"
	"errors"
	"time"

	bodyMeasurementDtos "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurement/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type bodyMeasurementDB struct {
	Id            uuid.UUID `gorm:"type:uuid;primaryKey"`
	ClientId uuid.UUID `gorm:"type:uuid;not null;index"`

	Mass          float64   `gorm:"type:decimal(10,2)"`
	Stature       float64   `gorm:"type:decimal(10,2)"`
	SittingHeight float64   `gorm:"column:sitting_height;type:decimal(10,2)"`
	ArmSpan       float64   `gorm:"column:arm_span;type:decimal(10,2)"`

	// Skinfolds
	Triceps      float64 `gorm:"type:decimal(10,2)"`
	Subscapular  float64 `gorm:"type:decimal(10,2)"`
	Biceps       float64 `gorm:"type:decimal(10,2)"`
	IliacCrest   float64 `gorm:"type:decimal(10,2)"`
	Supraspinale float64 `gorm:"type:decimal(10,2)"`
	Abdominal    float64 `gorm:"type:decimal(10,2)"`
	FrontThigh   float64 `gorm:"column:front_thigh;type:decimal(10,2)"`
	MedialCalf   float64 `gorm:"column:Medial_calf;type:decimal(10,2)"`

	// Girths/Circumferences
	Head       float64 `gorm:"type:decimal(10,2)"`
	Neck       float64 `gorm:"type:decimal(10,2)"`
	ArmRelaxed float64 `gorm:"column:arm_relaxed;type:decimal(10,2)"`
	ArmFlex    float64 `gorm:"column:arm_flex;type:decimal(10,2)"`
	Forearm    float64 `gorm:"type:decimal(10,2)"`
	Wrist      float64 `gorm:"type:decimal(10,2)"`
	Chest      float64 `gorm:"type:decimal(10,2)"`
	Waist      float64 `gorm:"type:decimal(10,2)"`
	Hip        float64 `gorm:"type:decimal(10,2)"`
	ThighHigh  float64 `gorm:"column:thigh_high;type:decimal(10,2)"`
	ThighLow   float64 `gorm:"column:thigh_low;type:decimal(10,2)"`
	Calf       float64 `gorm:"type:decimal(10,2)"`
	Ankle      float64 `gorm:"type:decimal(10,2)"`
	// Metadata & Relations
	CreatedAt time.Time `gorm:"column:created_at"`

}

func (e *bodyMeasurementDB) toEntity() (BodyMeasurement, *core.BaseError) {
	var errList []string

	validate := func(val float64, unit valueobject.MeasureUnit) valueobject.Measurer {
		m, err := valueobject.NewMeasurement(val, unit)
		if err != nil {
			errList = append(errList, err.Details...)
		}
		return m
	}

	// Basic Stats
	mass := validate(e.Mass, valueobject.KG)
	stature := validate(e.Stature, valueobject.CM)
	sittingHeight := validate(e.SittingHeight, valueobject.CM)
	armSpan := validate(e.ArmSpan, valueobject.CM)

	/* SkinFolds - Standardized to MM usually, change to CM if preferred */
	triceps := validate(e.Triceps, valueobject.MM)
	subscapular := validate(e.Subscapular, valueobject.MM)
	biceps := validate(e.Biceps, valueobject.MM)
	iliacCrest := validate(e.IliacCrest, valueobject.MM)
	supraspinale := validate(e.Supraspinale, valueobject.MM)
	abdominal := validate(e.Abdominal, valueobject.MM)
	frontThigh := validate(e.FrontThigh, valueobject.MM)
	medialCalf := validate(e.MedialCalf, valueobject.MM)

	/* Girths */
	head := validate(e.Head, valueobject.CM)
	neck := validate(e.Neck, valueobject.CM)
	armRelaxed := validate(e.ArmRelaxed, valueobject.CM)
	armFlex := validate(e.ArmFlex, valueobject.CM)
	forearm := validate(e.Forearm, valueobject.CM)
	wrist := validate(e.Wrist, valueobject.CM)
	chest := validate(e.Chest, valueobject.CM)
	waist := validate(e.Waist, valueobject.CM)
	hip := validate(e.Hip, valueobject.CM)
	thighHigh := validate(e.ThighHigh, valueobject.CM)
	thighLow := validate(e.ThighLow, valueobject.CM)
	calf := validate(e.Calf, valueobject.CM)
	ankle := validate(e.Ankle, valueobject.CM)

	createdAt, err := valueobject.NewTrackerFromTime(&e.CreatedAt)
	if err != nil {
		errList = append(errList, err.Details...)
	}

	clientId := valueobject.IdentifierFromValue(e.ClientId)

	if len(errList) > 0 {
		return nil, core.CorrupetedDatabase(errList)
	}

	return &entity{
		id:            valueobject.IdentifierFromValue(e.Id),
		clientId:      clientId,
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
		createdAt:     createdAt,
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

func (repository) TableName() string {
	return "body_measurements"
}

func NewRepository(db *gorm.DB) (Repository, *core.BaseError) {
	if err := db.AutoMigrate(&bodyMeasurementDB{}); err != nil {
		return nil, core.UnexpectedError(err.Error())
	}
	return &repository{db}, nil
}

func (r *repository) GetAllByUser(
	ctx context.Context,
	userID valueobject.Identifier,
) ([]BodyMeasurement, *core.BaseError) {
	var rawList []bodyMeasurementDB

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
		Model(&bodyMeasurementDB{}).
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
	var raw bodyMeasurementDB

	result := r.db.WithContext(ctx).
		First(&raw, "created_at::date = ?::date AND user_id = ?", date.ToTime(), userId.Key())
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
	var rawList []bodyMeasurementDB

	result := r.db.WithContext(ctx).
		Where("user_id = ?", clientId.Key()).
		Where("created_at:date BETWEEN ?::date AND ?::date", rangeDate.Start().ToTime(), rangeDate.End().ToTime()).
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
	var raw bodyMeasurementDB

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
	data := measurement.ToDB()
	result := r.db.WithContext(ctx).Create(&data)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id valueobject.Identifier) *core.BaseError {
	result := r.db.WithContext(ctx).Delete(&bodyMeasurementDB{}, id.Key())
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return BodyMeasurementNotFound(id.String())
	}

	return nil
}
