package interfaces

import (
	"context"

	"github.com/JooseMM/nutripia-backend-api/internal/users/models"
	"github.com/JooseMM/nutripia-backend-api/pkg/errors"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]*models.User, *errors.BaseError)
	GetById(ctx context.Context, id *uuid.UUID) (*models.User, *errors.BaseError)

	Create(ctx context.Context, u *models.User) *errors.BaseError
	Update(ctx context.Context, u *models.User) *errors.BaseError
	Delete(ctx context.Context, id *uuid.UUID) *errors.BaseError
}
