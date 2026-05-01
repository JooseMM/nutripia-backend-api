package bodyMeasurement

import (
	"context"

	bodyMeasurementDtos "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type BodyMeasurementManager interface {
	Create(
		ctx context.Context,
		dto bodyMeasurementDtos.CreateMeasurement,
		userId valueobject.Identifier,
	) (valueobject.Identifier, *core.BaseError)
	GetById(
		ctx context.Context,
		id valueobject.Identifier,
	) (BodyMeasurement, *core.BaseError)
	GetByRange(
		ctx context.Context,
		userId valueobject.Identifier,
		dateRange bodyMeasurementDtos.RangeDate,
	) ([]BodyMeasurement, *core.BaseError)
	DeleteOne(ctx context.Context, id valueobject.Identifier) *core.BaseError
}

type service struct {
	Repo Repository
}

func NewService(repo Repository) BodyMeasurementManager {
	return &service{repo}
}

func (u *service) Create(
	ctx context.Context,
	dto bodyMeasurementDtos.CreateMeasurement,
	userId valueobject.Identifier,
) (valueobject.Identifier, *core.BaseError) {
	now := valueobject.NewTracker()
	isAlreadyFill, err := u.Repo.DateAlreadyFill(ctx, now, userId)
	if err != nil {
		return nil, err
	}
	if *isAlreadyFill {
		return nil, BodyMeasurementAlreadyExisting(now.ToTime())
	}

	bm := NewBodyMeasurement(dto, userId)

	if err := u.Repo.Create(ctx, bm); err != nil {
		return nil, err
	}

	return bm.Id(), nil
}

func (u *service) GetById(
	ctx context.Context,
	id valueobject.Identifier,
) (BodyMeasurement, *core.BaseError) {
	bm, err := u.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return bm, nil
}

func (u *service) GetByRange(
	ctx context.Context,
	clientId valueobject.Identifier,
	rageDate bodyMeasurementDtos.RangeDate,
) ([]BodyMeasurement, *core.BaseError) {
	bm, err := u.Repo.GetByRange(ctx, rageDate, clientId)
	if err != nil {
		return nil, err
	}

	return bm, nil
}

func (u *service) DeleteOne(
	ctx context.Context,
	id valueobject.Identifier,
) *core.BaseError {
	err := u.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
