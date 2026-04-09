package nutritionist

import (
	"context"
	"time"

	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types"
	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type INutritionistService interface {
	GetById(id *uuid.UUID, ctx context.Context) (*nutritionistTypes.Nutritionist, *core.BaseError)
	DeleteOne(id *uuid.UUID, ctx context.Context) *core.BaseError
	UpdateIdentityInformation(
		id *uuid.UUID,
		userDto *nutritionistDtos.UpdateNutritionistIdentityRequest,
		ctx context.Context,
	) *core.BaseError
}

type NutritionistService struct {
	Repo INutritionistRepository
}

func NewUserService(repo INutritionistRepository) INutritionistService {
	return &NutritionistService{repo}
}


func (u *NutritionistService) GetById(
	id *uuid.UUID,
	ctx context.Context,
) (*nutritionistTypes.Nutritionist, *core.BaseError) {
	foundUser, unexpectedErr := u.Repo.GetById(ctx, id)
	if unexpectedErr != nil {
		return nil, unexpectedErr
	}

	return foundUser, nil
}

func (u *NutritionistService) DeleteOne(
	id *uuid.UUID,
	ctx context.Context,
) *core.BaseError {
	err := u.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *NutritionistService) UpdateIdentityInformation(
	id *uuid.UUID,
	userDto *nutritionistDtos.UpdateNutritionistIdentityRequest,
	ctx context.Context,
) *core.BaseError {
	foundUser, unexpectedErr := u.Repo.GetById(ctx, id)
	if unexpectedErr != nil {
		return unexpectedErr
	}

	foundEmailOwner, queryEmailErr := u.Repo.GetByEmailAddress(ctx, userDto.EmailAddress)
	if queryEmailErr != nil && queryEmailErr.ErrorCode != string(NUTRITIONIST_NOT_FOUND) {
		return queryEmailErr
	}
	if foundEmailOwner != nil && foundEmailOwner.ID != *id {
		return UserEmailAlreadyExisting(userDto.EmailAddress)
	}

	foundUser.Firstname = userDto.Firstname
	foundUser.Lastname = userDto.Lastname
	foundUser.EmailAddress = userDto.EmailAddress
	foundUser.DateBirth = userDto.BirthDate

	foundUser.UpdatedAt = time.Now()

	return u.Repo.Update(ctx, foundUser)
}
