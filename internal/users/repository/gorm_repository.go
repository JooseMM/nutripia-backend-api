package repository

import (
	"context"

	"github.com/JooseMM/nutripia-backend-api/internal/users/models"
	"github.com/JooseMM/nutripia-backend-api/internal/users/repository/interfaces"
	"github.com/JooseMM/nutripia-backend-api/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &postgresRepository{db}
}

func (r *postgresRepository) GetAll(ctx context.Context) ([]*models.User, *errors.BaseError) {
	var list []*models.User

	result := r.db.WithContext(ctx).Find(&list)
	if result.Error != nil {
		return nil, errors.UnexpectedError(result.Error.Error())
	}

	return list, nil
}

func (r *postgresRepository) GetById(
	ctx context.Context,
	id *uuid.UUID,
) (*models.User, *errors.BaseError) {
	var user models.User

	result := r.db.WithContext(ctx).First(&user, "ID = ?", id)
	if result.Error != nil {
		return nil, errors.UnexpectedError(result.Error.Error())
	}

	return &user, nil
}

func (r *postgresRepository) Create(ctx context.Context, u *models.User) *errors.BaseError {
	result := r.db.WithContext(ctx).Create(u)
	if result.Error != nil {
		return errors.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *postgresRepository) Update(ctx context.Context, u *models.User) *errors.BaseError {
	result := r.db.WithContext(ctx).Save(u)
	if result.Error != nil {
		return errors.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id *uuid.UUID) *errors.BaseError {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return errors.UnexpectedError(result.Error.Error())
	}

	return nil
}
