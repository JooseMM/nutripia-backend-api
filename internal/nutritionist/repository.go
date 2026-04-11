package nutritionist

import (
	"context"
	"errors"

	"github.com/JooseMM/nutripia-backend-api/internal/nutritionist/types"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NutritonistRepository struct {
	db *gorm.DB
}

type INutritionistRepository interface {
	GetAll(ctx *context.Context) ([]*nutritionistTypes.Nutritionist, *core.BaseError)
	GetById(ctx *context.Context, id *uuid.UUID) (*nutritionistTypes.Nutritionist, *core.BaseError)
	GetByEmailAddress(
		ctx *context.Context,
		emalAddress string,
	) (*nutritionistTypes.Nutritionist, *core.BaseError)

	Create(ctx *context.Context, u *nutritionistTypes.Nutritionist) *core.BaseError
	Update(ctx *context.Context, u *nutritionistTypes.Nutritionist) *core.BaseError
	Delete(ctx *context.Context, id *uuid.UUID) *core.BaseError
}

func NewNutritionistRepository(db *gorm.DB) INutritionistRepository {
	return &NutritonistRepository{db}
}

func (r *NutritonistRepository) GetAll(
	ctx *context.Context,
) ([]*nutritionistTypes.Nutritionist, *core.BaseError) {
	var list []*nutritionistTypes.Nutritionist

	result := r.db.WithContext(*ctx).Find(&list)
	if result.Error != nil {
		return nil, core.UnexpectedError(result.Error.Error())
	}

	return list, nil
}

func (r *NutritonistRepository) GetById(
	ctx *context.Context,
	id *uuid.UUID,
) (*nutritionistTypes.Nutritionist, *core.BaseError) {
	var user nutritionistTypes.Nutritionist

	result := r.db.WithContext(*ctx).First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, UserNotFound(id.String())
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	return &user, nil
}

func (r *NutritonistRepository) Create(
	ctx *context.Context,
	u *nutritionistTypes.Nutritionist,
) *core.BaseError {
	result := r.db.WithContext(*ctx).Create(u)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *NutritonistRepository) Update(
	ctx *context.Context,
	u *nutritionistTypes.Nutritionist,
) *core.BaseError {
	result := r.db.WithContext(*ctx).Save(u)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *NutritonistRepository) Delete(ctx *context.Context, id *uuid.UUID) *core.BaseError {
	result := r.db.WithContext(*ctx).Delete(&nutritionistTypes.Nutritionist{}, id)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return UserNotFound(id.String())
	}

	return nil
}

func (r *NutritonistRepository) GetByEmailAddress(
	ctx *context.Context,
	emailAddress string,
) (*nutritionistTypes.Nutritionist, *core.BaseError) {
	var user nutritionistTypes.Nutritionist

	result := r.db.WithContext(*ctx).First(&user, "email_address = ?", emailAddress)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, UserNotFound(emailAddress)
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	return &user, nil
}
