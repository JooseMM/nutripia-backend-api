package nutritionist

import (
	"context"

	nutritionistDtos "github.com/JooseMM/nutripia-backend-api/internal/nutritionist/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/JooseMM/nutripia-backend-api/pkg/core/valueobject"
)

type NutritionistManager interface {
	GetById(ctx context.Context, id valueobject.Identifier) (Nutritionist, *core.BaseError)
	Delete(ctx context.Context, id valueobject.Identifier) *core.BaseError
	Update(
		ctx context.Context,
		id valueobject.Identifier,
		data nutritionistDtos.UpdateNutritionist,
	) *core.BaseError
}

type NutritionistService struct {
	Repo Repository
}

func NewService(repo Repository) NutritionistManager {
	return &NutritionistService{repo}
}

func (u *NutritionistService) GetById(
	ctx context.Context,
	id valueobject.Identifier,
) (Nutritionist, *core.BaseError) {
	foundUser, unexpectedErr := u.Repo.GetById(ctx, id)
	if unexpectedErr != nil {
		return nil, unexpectedErr
	}

	return foundUser, nil
}

func (u *NutritionistService) Delete(
	ctx context.Context,
	id valueobject.Identifier,
) *core.BaseError {
	err := u.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *NutritionistService) Update(
	ctx context.Context,
	id valueobject.Identifier,
	data nutritionistDtos.UpdateNutritionist,
) *core.BaseError {
	foundUser, unexpectedErr := u.Repo.GetById(ctx, id)
	if unexpectedErr != nil {
		return unexpectedErr
	}

	foundEmailOwner, queryEmailErr := u.Repo.GetByEmailAddress(ctx, data.EmailAddress())
	if queryEmailErr != nil && queryEmailErr.ErrorCode != string(NUTRITIONIST_NOT_FOUND) {
		return queryEmailErr
	}

	if foundEmailOwner != nil && !foundEmailOwner.Id().IsEqual(id) {
		return UserEmailAlreadyExisting(data.EmailAddress().String())
	}

	foundUser.Update(data)
	return u.Repo.Update(ctx, foundUser)
}
