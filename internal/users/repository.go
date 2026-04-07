package users

import (
	"context"
	"errors"

	userModels "github.com/JooseMM/nutripia-backend-api/internal/users/models"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	GetAll(ctx context.Context) ([]*userModels.User, *core.BaseError)
	GetById(ctx context.Context, id *uuid.UUID) (*userModels.User, *core.BaseError)
	GetByEmailAddress(ctx context.Context, emalAddress string) (*userModels.User, *core.BaseError)

	Create(ctx context.Context, u *userModels.User) *core.BaseError
	Update(ctx context.Context, u *userModels.User) *core.BaseError
	Delete(ctx context.Context, id *uuid.UUID) *core.BaseError
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &postgresRepository{db}
}

func (r *postgresRepository) GetAll(ctx context.Context) ([]*userModels.User, *core.BaseError) {
	var list []*userModels.User

	result := r.db.WithContext(ctx).Find(&list)
	if result.Error != nil {
		return nil, core.UnexpectedError(result.Error.Error())
	}

	return list, nil
}

func (r *postgresRepository) GetById(
	ctx context.Context,
	id *uuid.UUID,
) (*userModels.User, *core.BaseError) {
	var user userModels.User

	result := r.db.WithContext(ctx).First(&user, "ID = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, NotFound(id.String())
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	return &user, nil
}

func (r *postgresRepository) Create(ctx context.Context, u *userModels.User) *core.BaseError {
	result := r.db.WithContext(ctx).Create(u)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *postgresRepository) Update(ctx context.Context, u *userModels.User) *core.BaseError {
	result := r.db.WithContext(ctx).Save(u)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id *uuid.UUID) *core.BaseError {
	result := r.db.WithContext(ctx).Delete(&userModels.User{}, id)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return NotFound(id.String())
	}

	return nil
}

func (r *postgresRepository) GetByEmailAddress(
	ctx context.Context,
	emailAddress string,
) (*userModels.User, *core.BaseError) {
	var user userModels.User

	result := r.db.WithContext(ctx).First(&user, "email_address = ?", emailAddress)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, NotFound(emailAddress)
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	return &user, nil
}
