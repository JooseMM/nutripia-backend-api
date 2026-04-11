package clients

import (
	"context"
	"errors"
	"fmt"

	"github.com/JooseMM/nutripia-backend-api/internal/clients/types"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClienRepository struct {
	db *gorm.DB
}

type IClientRepository interface {
	GetAllByNutritionist(
		userId *uuid.UUID,
		ctx context.Context,
	) ([]*clientTypes.Client, *core.BaseError)
	GetById(ctx context.Context, id *uuid.UUID) (*clientTypes.Client, *core.BaseError)
	GetByEmailAddress(
		ctx context.Context,
		emalAddress string,
	) (*clientTypes.Client, *core.BaseError)

	Create(ctx context.Context, u *clientTypes.Client) *core.BaseError
	Update(ctx context.Context, u *clientTypes.Client) *core.BaseError
	Delete(ctx context.Context, id *uuid.UUID) *core.BaseError
}

func NewClientRepository(db *gorm.DB) IClientRepository {
	return &ClienRepository{db}
}

func (r *ClienRepository) GetAllByNutritionist(
	userId *uuid.UUID,
	ctx context.Context,
) ([]*clientTypes.Client, *core.BaseError) {
	var list []*clientTypes.Client

	result := r.db.WithContext(ctx).Where("nutritionist_owner_id = ?", userId).Find(&list)
	if result.Error != nil {
		fmt.Printf("error: %s", result.Error)
		return nil, core.UnexpectedError(result.Error.Error())
	}

	return list, nil
}

func (r *ClienRepository) GetById(
	ctx context.Context,
	id *uuid.UUID,
) (*clientTypes.Client, *core.BaseError) {
	var user clientTypes.Client

	result := r.db.WithContext(ctx).First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, UserNotFound(id.String())
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	return &user, nil
}

func (r *ClienRepository) Create(ctx context.Context, u *clientTypes.Client) *core.BaseError {
	result := r.db.WithContext(ctx).Create(u)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *ClienRepository) Update(ctx context.Context, u *clientTypes.Client) *core.BaseError {
	result := r.db.WithContext(ctx).Save(u)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	return nil
}

func (r *ClienRepository) Delete(ctx context.Context, id *uuid.UUID) *core.BaseError {
	result := r.db.WithContext(ctx).Delete(&clientTypes.Client{}, id)
	if result.Error != nil {
		return core.UnexpectedError(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return UserNotFound(id.String())
	}

	return nil
}

func (r *ClienRepository) GetByEmailAddress(
	ctx context.Context,
	emailAddress string,
) (*clientTypes.Client, *core.BaseError) {
	var user clientTypes.Client

	result := r.db.WithContext(ctx).First(&user, "email_address = ?", emailAddress)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, UserNotFound(emailAddress)
		}
		return nil, core.UnexpectedError(result.Error.Error())
	}

	return &user, nil
}
