package interfaces

import (
	"context"

	"github.com/JooseMM/nutripia-backend-api/internal/users/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]*models.User, error)
	GetById(ctx context.Context, id *uuid.UUID) (*models.User, error)

	Create(ctx context.Context, u *models.User) error
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id *uuid.UUID) error
}
