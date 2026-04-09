package bodyMeasurement

import (
	"context"
	"errors"
	"time"

	bodyMeasurementTypes "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements/types"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BodyMeasurementRepository struct {
	db *gorm.DB
}

type IBodyMeasurementRepository interface {
	GetAllByUser(
		ctx context.Context,
		userId uuid.UUID,
	) ([]*bodyMeasurementTypes.BodyMeasurement, *core.BaseError)
	GetByDate(
		ctx context.Context,
		date *time.Time,
		userId *uuid.UUID,
	) (*bodyMeasurementTypes.BodyMeasurement, *core.BaseError)
	GetByRange(
		ctx context.Context,
		start *time.Time,
		end *time.Time,
		userId *uuid.UUID,
	) ([]*bodyMeasurementTypes.BodyMeasurement, *core.BaseError)
	GetById(
		ctx context.Context,
		id *uuid.UUID,
	) (*bodyMeasurementTypes.BodyMeasurement, *core.BaseError)
	Create(ctx context.Context, u *bodyMeasurementTypes.BodyMeasurement) *core.BaseError
	Delete(ctx context.Context, id *uuid.UUID) *core.BaseError
}

func NewBodyMeasurementRepository(db *gorm.DB) IBodyMeasurementRepository {
	return &BodyMeasurementRepository{db}
}

func (r *BodyMeasurementRepository) GetAllByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]*bodyMeasurementTypes.BodyMeasurement, *core.BaseError) {
	var list []*bodyMeasurementTypes.BodyMeasurement

	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&list)
	if result.Error != nil {
		return nil, core.UnexpectedError(result.Error.Error())
	}

	return list, nil
}

func (r *BodyMeasurementRepository) GetByDate(
	ctx context.Context,
	date *time.Time,
	userId *uuid.UUID,
) (*bodyMeasurementTypes.BodyMeasurement, *core.BaseError) {
	var measurement bodyMeasurementTypes.BodyMeasurement

	result := r.db.WithContext(ctx).First(&measurement, "created_at = ? AND user_id = ?", date, userId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, BodyMeasurementNotFound(date.Format("2000-01-01"))
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	return &measurement, nil
}

func (r *BodyMeasurementRepository) GetByRange(
	ctx context.Context,
	start *time.Time,
	end *time.Time,
	userId *uuid.UUID,
) ([]*bodyMeasurementTypes.BodyMeasurement, *core.BaseError) {
	var measurement []*bodyMeasurementTypes.BodyMeasurement

	result := r.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at ASC").
		Find(&measurement)

	if result.Error != nil {
		return nil, core.UnexpectedError(result.Error.Error())
	}
	return measurement, nil
}

func (r *BodyMeasurementRepository) GetById(
	ctx context.Context,
	id *uuid.UUID,
) (*bodyMeasurementTypes.BodyMeasurement, *core.BaseError) {
	var measurement bodyMeasurementTypes.BodyMeasurement

	result := r.db.WithContext(ctx).First(&measurement, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, BodyMeasurementNotFound(id.String())
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	return &measurement, nil
}

func (r *BodyMeasurementRepository) Create(
	ctx context.Context,
	measurement *bodyMeasurementTypes.BodyMeasurement,
) *core.BaseError {
	result := r.db.WithContext(ctx).Create(measurement)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *BodyMeasurementRepository) Delete(ctx context.Context, id *uuid.UUID) *core.BaseError {
	result := r.db.WithContext(ctx).Delete(&bodyMeasurementTypes.BodyMeasurement{}, id)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return BodyMeasurementNotFound(id.String())
	}

	return nil
}
