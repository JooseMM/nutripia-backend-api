package bodyMeasurement

import (
	"context"
	"time"

	bodyMeasurementTypes "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements/types"
	bodyMeasurementDtos "github.com/JooseMM/nutripia-backend-api/internal/bodyMeasurements/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type IBodyMeasurementService interface {
	Create(
		dto *bodyMeasurementDtos.CreateMeasurementDto,
		ctx context.Context,
	) (*uuid.UUID, *core.BaseError)
	GetById(
		id *uuid.UUID,
		ctx context.Context,
	) (*bodyMeasurementTypes.BodyMeasurement, *core.BaseError)
	GetByRange(
		userId *uuid.UUID,
		start *time.Time,
		end *time.Time,
		ctx context.Context,
	) ([]*bodyMeasurementTypes.BodyMeasurement, *core.BaseError)
	DeleteOne(id *uuid.UUID, ctx context.Context) *core.BaseError
}

type BodyMeasurementService struct {
	Repo IBodyMeasurementRepository
}

func NewBodyMeasurementService(repo IBodyMeasurementRepository) IBodyMeasurementService {
	return &BodyMeasurementService{repo}
}

func (u *BodyMeasurementService) Create(
	dto *bodyMeasurementDtos.CreateMeasurementDto,
	ctx context.Context,
) (*uuid.UUID, *core.BaseError) {
	now := time.Now()
	query, queryErr := u.Repo.GetByDate(ctx, &now, &dto.ClientId)
	if queryErr != nil && queryErr.ErrorCode != string(BODY_MEASUREMENT_RECORD_NOT_FOUND) {
		return nil, queryErr
	}
	if query != nil {
		return nil, BodyMeasurementAlreadyExisting(now)
	}

	measurement := &bodyMeasurementTypes.BodyMeasurement{
		ID:            uuid.New(),
		Mass:          dto.Mass,
		Stature:       dto.Stature,
		SittingHeight: dto.SittingHeight,
		ArmSpan:       dto.ArmSpan,
		Triceps:       dto.Triceps,
		Subscapular:   dto.Subscapular,
		Biceps:        dto.Biceps,
		IliacCrest:    dto.IliacCrest,
		Supraspinale:  dto.Supraspinale,
		Abdominal:     dto.Abdominal,
		FrontThigh:    dto.FrontThigh,
		MedialCalf:    dto.MedialCalf,
		Head:          dto.Head,
		Neck:          dto.Neck,
		ArmRelaxed:    dto.ArmRelaxed,
		ArmFlex:       dto.ArmFlex,
		Forearm:       dto.Forearm,
		Wrist:         dto.Wrist,
		Chest:         dto.Chest,
		Waist:         dto.Waist,
		Hip:           dto.Hip,
		ThighHigh:     dto.ThighHigh,
		ThighLow:      dto.ThighLow,
		Calf:          dto.Calf,
		Ankle:         dto.Ankle,

		ClientId:    dto.ClientId,
		CreatedAt: time.Now(),
	}

	if err := u.Repo.Create(ctx, measurement); err != nil {
		return nil, err
	}

	return &measurement.ID, nil
}

func (u *BodyMeasurementService) GetById(
	id *uuid.UUID,
	ctx context.Context,
) (*bodyMeasurementTypes.BodyMeasurement, *core.BaseError) {
	foundUser, unexpectedErr := u.Repo.GetById(ctx, id)
	if unexpectedErr != nil {
		return nil, unexpectedErr
	}

	return foundUser, nil
}

func (u *BodyMeasurementService) GetByRange(
	userId *uuid.UUID,
	start *time.Time,
	end *time.Time,
	ctx context.Context,
) ([]*bodyMeasurementTypes.BodyMeasurement, *core.BaseError) {
	found, unexpectedErr := u.Repo.GetByRange(ctx, start, end, userId)
	if unexpectedErr != nil {
		return nil, unexpectedErr
	}

	return found, nil
}

func (u *BodyMeasurementService) DeleteOne(
	id *uuid.UUID,
	ctx context.Context,
) *core.BaseError {
	err := u.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
