package users

import (
	"context"
	userModels "github.com/JooseMM/nutripia-backend-api/internal/users/models"
	userObjects "github.com/JooseMM/nutripia-backend-api/internal/users/models/value_objects"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
	"time"
)

type IUserService interface {
	CreateUser(dto userModels.CreateUserDto, ctx context.Context) (*uuid.UUID, *core.BaseError)
}

type UserService struct {
	Repo IUserRepository
}

func NewUserService(repo IUserRepository) IUserService {
	return &UserService{repo}
}

func (u *UserService) CreateUser(
	dto userModels.CreateUserDto,
	ctx context.Context,
) (*uuid.UUID, *core.BaseError) {

	foundUser, unexpectedErr := u.Repo.GetByEmailAddress(ctx, dto.EmailAddress)
	if foundUser != nil {
		return nil, EmailAlreadyExisting(dto.EmailAddress)
	}

	now := time.Now()
	user := &userModels.User{
		UserIdentity: userObjects.UserIdentity{
			ID:           uuid.New(),
			Firstname:    dto.Firstname,
			Lastname:     dto.Lastname,
			EmailAddress: dto.EmailAddress,
			DateBirth:    dto.BirthDate,
		},
		AuthenticationInformation: userObjects.AuthenticationInformation{},
		TrackingInformation: userObjects.TrackingInformation{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	err := u.Repo.Create(ctx, user)
	if err != nil {
		return nil, unexpectedErr
	}

	return &user.ID, nil
}
