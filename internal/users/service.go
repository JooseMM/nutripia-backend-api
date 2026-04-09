package users

import (
	"context"
	"time"

	"github.com/JooseMM/nutripia-backend-api/internal/security"
	"github.com/JooseMM/nutripia-backend-api/internal/users/types"
	"github.com/JooseMM/nutripia-backend-api/internal/users/types/dtos"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
	"github.com/google/uuid"
)

type IUserService interface {
	CreateUser(dto *userDtos.CreateUserRequest, role userTypes.Role, ctx context.Context) (*uuid.UUID, *core.BaseError)
	GetById(id *uuid.UUID, ctx context.Context) (*userTypes.User, *core.BaseError)
	DeleteOne(id *uuid.UUID, ctx context.Context) *core.BaseError
	UpdateIdentityInformation(
		id *uuid.UUID,
		userDto *userDtos.UpdateIdentityRequest,
		ctx context.Context,
	) *core.BaseError
}

type UserService struct {
	Repo IUserRepository
}

func NewUserService(repo IUserRepository) IUserService {
	return &UserService{repo}
}

func (u *UserService) CreateUser(
	userDto *userDtos.CreateUserRequest,
	role userTypes.Role,
	ctx context.Context,
) (*uuid.UUID, *core.BaseError) {

	foundEmailOwner, queryEmailErr := u.Repo.GetByEmailAddress(ctx, userDto.EmailAddress)
	if queryEmailErr != nil && queryEmailErr.ErrorCode != string(NOT_FOUND) {
		return nil, queryEmailErr
	}
	if foundEmailOwner != nil {
		return nil, UserEmailAlreadyExisting(userDto.EmailAddress)
	}

	hash, hashErr := security.HashPassword(userDto.Password)
	if hashErr != nil {
		return nil, core.UnexpectedError(hashErr.Error())
	}

	now := time.Now()
	user := &userTypes.User{
		ID: uuid.New(),
		UserIdentity: userTypes.UserIdentity{
			Firstname:    userDto.Firstname,
			Lastname:     userDto.Lastname,
			EmailAddress: userDto.EmailAddress,
			DateBirth:    userDto.BirthDate,
		},
		AuthenticationInformation: userTypes.AuthenticationInformation{
			Role:             role,
			PasswordHash:     hash,
			IsEmailConfirmed: false,
		},
		TrackingInformation: userTypes.TrackingInformation{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	creationErr := u.Repo.Create(ctx, user)
	if creationErr != nil {
		return nil, creationErr
	}

	return &user.ID, nil
}

func (u *UserService) GetById(
	id *uuid.UUID,
	ctx context.Context,
) (*userTypes.User, *core.BaseError) {

	foundUser, unexpectedErr := u.Repo.GetById(ctx, id)
	if unexpectedErr != nil {
		return nil, unexpectedErr
	}

	return foundUser, nil
}

func (u *UserService) DeleteOne(
	id *uuid.UUID,
	ctx context.Context,
) *core.BaseError {
	err := u.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) UpdateIdentityInformation(
	id *uuid.UUID,
	userDto *userDtos.UpdateIdentityRequest,
	ctx context.Context,
) *core.BaseError {
	foundUser, unexpectedErr := u.Repo.GetById(ctx, id)
	if unexpectedErr != nil {
		return unexpectedErr
	}

	foundEmailOwner, queryEmailErr := u.Repo.GetByEmailAddress(ctx, userDto.EmailAddress)
	if queryEmailErr != nil && queryEmailErr.ErrorCode != string(NOT_FOUND) {
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
