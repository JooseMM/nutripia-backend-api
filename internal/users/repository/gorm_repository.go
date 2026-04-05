package repository

import (
	"context"

	"github.com/JooseMM/nutripia-backend-api/internal/users/models"
	"github.com/JooseMM/nutripia-backend-api/internal/users/repository/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &postgresRepository{db}
}

func (r *postgresRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	var list []*models.User

	result := r.db.WithContext(ctx).Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}

	return list, nil
}

func (r *postgresRepository) GetById(ctx context.Context, id *uuid.UUID) (*models.User, error) {
	var user models.User

	result := r.db.WithContext(ctx).First(&user, "ID = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *postgresRepository) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *postgresRepository) Update(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *postgresRepository) Delete(ctx context.Context, id *uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(id).Error
}
